package clickhouse

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/hoppermq/streamly/pkg/domain"
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

func (t *Translator) Translate(ast *domain.QueryAstRequest) (*QueryBuilder, error) {
	builder := NewQueryBuilder()

	if err := t.translateSelect(ast.Select, builder); err != nil {
		return nil, fmt.Errorf("failed to translate SELECT: %w", err)
	}

	if err := t.translateFrom(ast.From, builder); err != nil {
		return nil, fmt.Errorf("failed to translate FROM: %w", err)
	}

	if err := t.translateWhere(ast.TimeRange, ast.Where, builder); err != nil {
		return nil, fmt.Errorf("failed to translate WHERE: %w", err)
	}

	if err := t.translateGroupBy(ast.GroupBy, builder); err != nil {
		return nil, fmt.Errorf("failed to translate GROUP BY: %w", err)
	}

	if err := t.translateOrderBy(ast.OrderBy, builder); err != nil {
		return nil, fmt.Errorf("failed to translate ORDER BY: %w", err)
	}

	if ast.Limit != nil {
		builder.SetLimit(*ast.Limit)
	}

	if ast.Offset != nil {
		builder.SetOffset(*ast.Offset)
	}

	return builder, nil
}

func (t *Translator) translateSelect(selectClauses []domain.SelectClause, builder *QueryBuilder) error {
	if len(selectClauses) == 0 {
		return fmt.Errorf("SELECT clause cannot be empty")
	}

	for _, clause := range selectClauses {
		if clause.IsField() {
			builder.SelectFields(*clause.Field)
		} else if clause.IsFunction() {
			fnArgs := strings.Join(clause.Function.Args, ", ")
			builder.SelectFunc(
				clause.Function.Function,
				fnArgs,
				clause.Function.Alias,
			)
		} else {
			return fmt.Errorf("unknown SELECT clause type")
		}
	}

	return nil
}

func (t *Translator) translateFrom(datasource domain.Datasource, builder *QueryBuilder) error {
	if datasource == "" {
		return fmt.Errorf("FROM datasource cannot be empty")
	}

	builder.From(string(datasource))
	return nil
}

func (t *Translator) translateWhere(timeRange domain.TimeRange, whereClauses []domain.WhereClause, builder *QueryBuilder) error {
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
				valueSlice, ok := where.Value.([]any)
				if ok {
					values = valueSlice
				} else {
					return fmt.Errorf("IN operator requires array value for field %s", where.Field)
				}
			}
			builder.WhereIn(where.Field, values)
		} else {
			builder.Where(where.Field, where.Op, where.Value)
		}
	}

	return nil
}

func (t *Translator) translateGroupBy(groupByClauses []domain.GroupByClause, builder *QueryBuilder) error {
	for _, gb := range groupByClauses {
		if gb.IsField() {
			builder.GroupBy(*gb.Field)
		} else if gb.IsTimeWindow() {
			field := gb.TimeWindow.Field
			if field == "" {
				field = "timestamp"
			}
			builder.GroupByTimeWindow(gb.TimeWindow.Window, field)
		} else {
			return fmt.Errorf("unknown GROUP BY clause type")
		}
	}

	return nil
}

func (t *Translator) translateOrderBy(orderByClauses []domain.OrderByClause, builder *QueryBuilder) error {
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
