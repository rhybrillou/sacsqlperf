package main

import (
	"fmt"
	"strings"
)

type qualifiedColumn struct {
	tableName  string
	columnName string
	value      interface{}
}

type innerJoin struct {
	left  qualifiedColumn
	right qualifiedColumn
}

type orderColumn struct {
	column   qualifiedColumn
	reversed bool
}

type pagination struct {
	limit  int
	offset int
}

type query struct {
	statement            string
	statementTargets     []string
	targetTables         []string
	innerJoins           []innerJoin
	whereClause          whereClausePart
	orderBy              []orderColumn
	queryPagination      *pagination
	scopeLevel           string
	scopeTable           string
	scopeClusterColumn   string
	scopeNamespaceColumn string
}

func (q *query) ForExecution() (string, []interface{}) {
	params := make([]interface{}, 0)
	var qb strings.Builder
	qb.WriteString(q.statement)
	if len(q.statementTargets) > 0 {
		qb.WriteString(" ")
		qb.WriteString(strings.Join(q.statementTargets, ", "))
	}
	qb.WriteString(" from ")
	qb.WriteString(strings.Join(q.targetTables, " "))
	for _, join := range q.innerJoins {
		qb.WriteString(
			fmt.Sprintf(
				" inner join %s on %s.%s = %s.%s",
				join.right.tableName,
				join.left.tableName,
				join.left.columnName,
				join.right.tableName,
				join.right.columnName,
			),
		)
	}
	if q.whereClause != nil {
		whereClause, bindValues := q.whereClause.asWhereClausePart()
		qb.WriteString(" where ")
		qb.WriteString(whereClause)
		params = append(params, bindValues...)
	}
	if len(q.orderBy) > 0 {
		qb.WriteString(" order by ")
		for ix, order := range q.orderBy {
			if ix > 0 {
				qb.WriteString(", ")
			}
			qb.WriteString(fmt.Sprintf("%s.%s", order.column.tableName, order.column.columnName))
			if order.reversed {
				qb.WriteString(" desc")
			}
		}
	}
	if q.queryPagination != nil {
		if q.queryPagination.offset > 0 {
			qb.WriteString(fmt.Sprintf(" offset %d", q.queryPagination.offset))
		}
		if q.queryPagination.limit > 0 {
			qb.WriteString(fmt.Sprintf(" limit %d", q.queryPagination.limit))
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
