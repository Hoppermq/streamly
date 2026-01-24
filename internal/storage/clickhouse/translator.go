package clickhouse

import (
	"log/slog"
	"strings"

	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/errors"
)

type Translator struct {
	logger *slog.Logger
}

// TODO: implement the time range.

type TranslatorOption func(*Translator)

func TranslatorWithLogger(logger *slog.Logger) TranslatorOption {
	return func(translator *Translator) {
		translator.logger = logger
	}
}

func (t *Translator) Translate(
	ast *domain.QueryAstRequest,
) (*QueryBuilder, error) {
	builder := NewQueryBuilder()

	if err := t.translateSelect(ast.Select, builder); err != nil {
		return nil, errors.TranslatorFailedToTranslate(
			errors.ErrSelectTranslationFailed,
			err,
		)
	}

	if err := t.translateFrom(ast.From, builder); err != nil {
		return nil, errors.TranslatorFailedToTranslate(
			errors.ErrFromTranslationFailed,
			err,
		)
	}

	if err := t.translateWhere(ast.TimeRange, ast.Where, builder); err != nil {
		return nil, errors.TranslatorFailedToTranslate(
			errors.ErrWhereTranslationFailed,
			err,
		)
	}

	if err := t.translateGroupBy(ast.GroupBy, builder); err != nil {
		return nil, errors.TranslatorFailedToTranslate(
			errors.ErrGroupByTranslationFailed,
			err,
		)
	}

	if err := t.translateOrderBy(ast.OrderBy, builder); err != nil {
		return nil, errors.TranslatorFailedToTranslate(
			errors.ErrOrderByTranslationFailed,
			err,
		)
	}

	if ast.Limit != nil {
		builder.SetLimit(*ast.Limit)
	}

	if ast.Offset != nil {
		builder.SetOffset(*ast.Offset)
	}

	return builder, nil
}

func (t *Translator) translateSelect(
	selectClauses []domain.SelectClause,
	builder *QueryBuilder,
) error {
	if len(selectClauses) == 0 {
		return errors.ErrSelectClauseEmpty
	}

	for _, clause := range selectClauses {
		switch {
		case clause.IsField():
			builder.SelectFields(*clause.Field)
		case clause.IsFunction():
			fnArgs := strings.Join(clause.Function.Args, ", ")
			builder.SelectFunc(
				clause.Function.Function,
				fnArgs,
				clause.Function.Alias,
			)
		default:
			return errors.ErrSelectClauseType
		}
	}

	return nil
}

func (t *Translator) translateFrom(
	datasource domain.Datasource,
	builder *QueryBuilder,
) error {
	if datasource == "" {
		return errors.ErrFromEmpty
	}

	builder.From(string(datasource))
	return nil
}

func (t *Translator) translateWhere(
	timeRange domain.TimeRange,
	whereClauses []domain.WhereClause,
	builder *QueryBuilder,
) error {
	if timeRange.Start != "" {
		builder.Where("timestamp", ">=", timeRange.Start)
	}
	if timeRange.End != "" {
		builder.Where("timestamp", "<=", timeRange.End)
	}

	for _, where := range whereClauses {
		if where.Op == "IN" {
			values, ok := where.Value.([]any)
			if !ok {
				return errors.TranslatorInOperatorInvalidValue(where.Field)
			}
			builder.WhereIn(where.Field, values)
		} else {
			builder.Where(where.Field, where.Op, where.Value)
		}
	}

	return nil
}

func (t *Translator) translateGroupBy(
	groupByClauses []domain.GroupByClause,
	builder *QueryBuilder,
) error {
	for _, gb := range groupByClauses {
		switch {
		case gb.IsField():
			builder.GroupBy(*gb.Field)
		case gb.IsTimeWindow():
			field := gb.TimeWindow.Field
			if field == "" {
				field = "timestamp"
			}
			builder.GroupByTimeWindow(gb.TimeWindow.Window, field)
		default:
			return errors.ErrUnknownGroupBy
		}
	}

	return nil
}

func (t *Translator) translateOrderBy(
	orderByClauses []domain.OrderByClause,
	builder *QueryBuilder,
) error {
	if orderByClauses == nil {
		return nil
	}
	for _, ob := range orderByClauses {
		direction := ob.Direction
		if direction == "" {
			direction = "DESC"
		}
		builder.OrderBy(ob.Field, direction)
	}

	return nil
}

func NewTranslator(opts ...TranslatorOption) *Translator {
	t := &Translator{}
	for _, opt := range opts {
		opt(t)
	}

	return t
}
