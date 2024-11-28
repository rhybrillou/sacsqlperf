package query

import (
	"fmt"
	"strings"
)

type QualifiedColumn struct {
	TableName  string
	ColumnName string
	Value      interface{}
}

type InnerJoin struct {
	Left  QualifiedColumn
	Right QualifiedColumn
}

type OrderColumn struct {
	Column   QualifiedColumn
	Reversed bool
}

type Pagination struct {
	Limit  int
	Offset int
}

type Query struct {
	Statement            string
	StatementTargets     []string
	TargetTables         []string
	InnerJoins           []InnerJoin
	WhereClause          WhereClausePart
	OrderBy              []OrderColumn
	QueryPagination      *Pagination
	ScopeLevel           string
	ScopeTable           string
	ScopeClusterColumn   string
	ScopeNamespaceColumn string
}

func (q *Query) ForExecution() (string, []interface{}) {
	params := make([]interface{}, 0)
	var qb strings.Builder
	qb.WriteString(q.Statement)
	if len(q.StatementTargets) > 0 {
		qb.WriteString(" ")
		qb.WriteString(strings.Join(q.StatementTargets, ", "))
	}
	qb.WriteString(" from ")
	qb.WriteString(strings.Join(q.TargetTables, " "))
	for _, join := range q.InnerJoins {
		qb.WriteString(
			fmt.Sprintf(
				" inner join %s on %s.%s = %s.%s",
				join.Right.TableName,
				join.Left.TableName,
				join.Left.ColumnName,
				join.Right.TableName,
				join.Right.ColumnName,
			),
		)
	}
	if q.WhereClause != nil {
		whereClause, bindValues := q.WhereClause.AsWhereClausePart()
		qb.WriteString(" where ")
		qb.WriteString(whereClause)
		params = append(params, bindValues...)
	}
	if len(q.OrderBy) > 0 {
		qb.WriteString(" order by ")
		for ix, order := range q.OrderBy {
			if ix > 0 {
				qb.WriteString(", ")
			}
			qb.WriteString(fmt.Sprintf("%s.%s", order.Column.TableName, order.Column.ColumnName))
			if order.Reversed {
				qb.WriteString(" desc")
			}
		}
	}
	if q.QueryPagination != nil {
		if q.QueryPagination.Offset > 0 {
			qb.WriteString(fmt.Sprintf(" offset %d", q.QueryPagination.Offset))
		}
		if q.QueryPagination.Limit > 0 {
			qb.WriteString(fmt.Sprintf(" limit %d", q.QueryPagination.Limit))
		}
	}
	return enumerateBindValues(qb.String()), params
}

func enumerateBindValues(statement string) string {
	parts := strings.Split(statement, "$$")
	var result strings.Builder
	for ix, part := range parts {
		if ix > 0 {
			result.WriteString(fmt.Sprintf("$%d", ix))
		}
		result.WriteString(part)
	}
	return result.String()
}
