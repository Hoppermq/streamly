package clickhouse

import (
	"fmt"
	"strings"

	"github.com/hoppermq/streamly/pkg/domain/errors"
)

type SelectExpr struct {
	Expression string
	Alias      string
}

type WhereExpr struct {
	Field    string
	Operator string
	Value    any
}

type GroupByExpr struct {
	Expression string
}

type OrderByExpr struct {
	Field     string
	Direction string
}

type SQLComponents struct {
	SelectClauses  []SelectExpr
	FromSource     string
	WhereClauses   []WhereExpr
	GroupByClauses []GroupByExpr
	OrderByClauses []OrderByExpr
	Limit          *int
	Offset         *int
}

type QueryBuilder struct {
	components SQLComponents
}

var allowedOperators = map[string]bool{
	"=": true, "!=": true, ">": true, "<": true,
	">=": true, "<=": true, "IN": true, "LIKE": true,
}

func (qb *QueryBuilder) Select(exprs ...SelectExpr) *QueryBuilder {
	qb.components.SelectClauses = append(qb.components.SelectClauses, exprs...)
	return qb
}

func (qb *QueryBuilder) SelectFields(fields ...string) *QueryBuilder {
	for _, field := range fields {
		qb.components.SelectClauses = append(qb.components.SelectClauses,
			SelectExpr{Expression: field})
	}
	return qb
}

func (qb *QueryBuilder) SelectAs(field, alias string) *QueryBuilder {
	qb.components.SelectClauses = append(qb.components.SelectClauses,
		SelectExpr{Expression: field, Alias: alias})
	return qb
}

func (qb *QueryBuilder) SelectFunc(fn, args, alias string) *QueryBuilder {
	expr := fmt.Sprintf("%s(%s)", fn, args)
	qb.components.SelectClauses = append(qb.components.SelectClauses,
		SelectExpr{Expression: expr, Alias: alias})
	return qb
}

func (qb *QueryBuilder) From(dataSource string) *QueryBuilder {
	qb.components.FromSource = dataSource
	return qb
}

func (qb *QueryBuilder) Where(field, operator string, value any) *QueryBuilder {
	if !allowedOperators[operator] {
		return qb
	}

	qb.components.WhereClauses = append(qb.components.WhereClauses,
		WhereExpr{
			Field:    field,
			Operator: operator,
			Value:    value,
		})
	return qb
}

func (qb *QueryBuilder) WhereIn(field string, values []any) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}

	qb.components.WhereClauses = append(qb.components.WhereClauses,
		WhereExpr{
			Field:    field,
			Operator: "IN",
			Value:    values,
		})
	return qb
}

func (qb *QueryBuilder) GroupBy(expressions ...string) *QueryBuilder {
	for _, expr := range expressions {
		qb.components.GroupByClauses = append(qb.components.GroupByClauses,
			GroupByExpr{Expression: expr})
	}
	return qb
}

func (qb *QueryBuilder) GroupByTimeWindow(window, field string) *QueryBuilder {
	expr := fmt.Sprintf("toStartOfInterval(%s, INTERVAL %s)", field, window)
	qb.components.GroupByClauses = append(qb.components.GroupByClauses,
		GroupByExpr{Expression: expr})
	return qb
}

func (qb *QueryBuilder) OrderBy(field, direction string) *QueryBuilder {
	if direction == "" {
		direction = "DESC"
	}

	qb.components.OrderByClauses = append(qb.components.OrderByClauses,
		OrderByExpr{
			Field:     field,
			Direction: strings.ToUpper(direction),
		})
	return qb
}

func (qb *QueryBuilder) SetLimit(limit int) *QueryBuilder {
	qb.components.Limit = &limit
	return qb
}

func (qb *QueryBuilder) SetOffset(offset int) *QueryBuilder {
	qb.components.Offset = &offset
	return qb
}

func (qb *QueryBuilder) Build() (string, []any, error) {
	if err := qb.validate(); err != nil {
		return "", nil, err
	}

	var sql strings.Builder
	var args []any

	qb.buildSelect(&sql)
	qb.buildFrom(&sql)

	whereArgs, err := qb.buildWhere(&sql)
	if err != nil {
		return "", nil, err
	}
	args = append(args, whereArgs...)

	qb.buildGroupBy(&sql)
	qb.buildOrderBy(&sql)

	limitOffsetArgs := qb.buildLimitOffset(&sql)
	args = append(args, limitOffsetArgs...)

	return sql.String(), args, nil
}

func (qb *QueryBuilder) validate() error {
	if len(qb.components.SelectClauses) == 0 {
		return errors.ErrNoSelectClauseDefined
	}
	if qb.components.FromSource == "" {
		return errors.ErrNoFromSourceDefined
	}
	return nil
}

func (qb *QueryBuilder) buildSelect(sql *strings.Builder) {
	sql.WriteString("SELECT ")
	for i := range qb.components.SelectClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(qb.components.SelectClauses[i].Expression)
		if qb.components.SelectClauses[i].Alias != "" {
			sql.WriteString(" AS ")
			sql.WriteString(qb.components.SelectClauses[i].Alias)
		}
	}
}

func (qb *QueryBuilder) buildFrom(sql *strings.Builder) {
	sql.WriteString(" FROM ")
	sql.WriteString(qb.components.FromSource)
}

func (qb *QueryBuilder) buildWhere(sql *strings.Builder) ([]any, error) {
	if len(qb.components.WhereClauses) == 0 {
		return nil, nil
	}

	var args []any
	sql.WriteString(" WHERE ")

	for i := range qb.components.WhereClauses {
		if i > 0 {
			sql.WriteString(" AND ")
		}

		whereArgs, err := qb.buildWhereClause(sql, qb.components.WhereClauses[i])
		if err != nil {
			return nil, err
		}
		args = append(args, whereArgs...)
	}

	return args, nil
}

func (qb *QueryBuilder) buildWhereClause(sql *strings.Builder, where WhereExpr) ([]any, error) {
	if where.Operator == "IN" {
		return qb.buildInClause(sql, where)
	}

	if _, err := fmt.Fprintf(sql, "%s %s ?", where.Field, where.Operator); err != nil {
		return nil, err
	}

	return []any{where.Value}, nil
}

func (qb *QueryBuilder) buildInClause(sql *strings.Builder, where WhereExpr) ([]any, error) {
	values, ok := where.Value.([]any)
	if !ok {
		return nil, errors.ErrInOperator
	}

	placeholders := make([]string, len(values))
	for j := range placeholders {
		placeholders[j] = "?"
	}

	_, err := fmt.Fprintf(
		sql,
		"%s IN (%s)",
		where.Field,
		strings.Join(placeholders, ", "),
	)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (qb *QueryBuilder) buildGroupBy(sql *strings.Builder) {
	if len(qb.components.GroupByClauses) == 0 {
		return
	}

	sql.WriteString(" GROUP BY ")
	for i, gb := range qb.components.GroupByClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(gb.Expression)
	}
}

func (qb *QueryBuilder) buildOrderBy(sql *strings.Builder) {
	if len(qb.components.OrderByClauses) == 0 {
		return
	}

	sql.WriteString(" ORDER BY ")
	for i := range qb.components.OrderByClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		_, err := fmt.Fprintf(
			sql,
			"%s %s",
			qb.components.OrderByClauses[i].Field,
			qb.components.OrderByClauses[i].Direction,
		)
		if err != nil {
			return
		}
	}
}

func (qb *QueryBuilder) buildLimitOffset(sql *strings.Builder) []any {
	var args []any

	if qb.components.Limit != nil {
		sql.WriteString(" LIMIT ?")
		args = append(args, *qb.components.Limit)
	}

	if qb.components.Offset != nil {
		sql.WriteString(" OFFSET ?")
		args = append(args, *qb.components.Offset)
	}

	return args
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}
