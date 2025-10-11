package clickhouse

import (
	"fmt"
	"strings"
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
	if len(qb.components.SelectClauses) == 0 {
		return "", nil, fmt.Errorf("no SELECT clauses defined")
	}

	if qb.components.FromSource == "" {
		return "", nil, fmt.Errorf("no FROM source defined")
	}

	var sql strings.Builder
	var args []any

	sql.WriteString("SELECT ")
	for i, sel := range qb.components.SelectClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(sel.Expression)
		if sel.Alias != "" {
			sql.WriteString(" AS ")
			sql.WriteString(sel.Alias)
		}
	}

	sql.WriteString(" FROM ")
	sql.WriteString(qb.components.FromSource)

	if len(qb.components.WhereClauses) > 0 {
		sql.WriteString(" WHERE ")
		for i, where := range qb.components.WhereClauses {
			if i > 0 {
				sql.WriteString(" AND ")
			}

			if where.Operator == "IN" {
				values, ok := where.Value.([]any)
				if !ok {
					return "", nil, fmt.Errorf("IN operator requires []any value")
				}

				placeholders := make([]string, len(values))
				for j := range placeholders {
					placeholders[j] = "?"
					args = append(args, values[j])
				}
				sql.WriteString(fmt.Sprintf("%s IN (%s)", where.Field, strings.Join(placeholders, ", ")))
			} else {
				sql.WriteString(fmt.Sprintf("%s %s ?", where.Field, where.Operator))
				args = append(args, where.Value)
			}
		}
	}

	if len(qb.components.GroupByClauses) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, gb := range qb.components.GroupByClauses {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(gb.Expression)
		}
	}

	if len(qb.components.OrderByClauses) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, ob := range qb.components.OrderByClauses {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(fmt.Sprintf("%s %s", ob.Field, ob.Direction))
		}
	}

	if qb.components.Limit != nil {
		sql.WriteString(" LIMIT ?")
		args = append(args, *qb.components.Limit)
	}

	if qb.components.Offset != nil {
		sql.WriteString(" OFFSET ?")
		args = append(args, *qb.components.Offset)
	}

	return sql.String(), args, nil
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}
