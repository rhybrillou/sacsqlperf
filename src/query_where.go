package main

import (
	"fmt"
	"strings"
)

type whereClausePart interface {
	asWhereClausePart() (string, []interface{})
}

func (qc *qualifiedColumn) asWhereClausePart() (string, []interface{}) {
	return fmt.Sprintf("%s.%s = $$", qc.tableName, qc.columnName), []interface{}{qc.value}
}

type wcOr struct {
	operands []whereClausePart
}

func (wcp *wcOr) asWhereClausePart() (string, []interface{}) {
	var qb strings.Builder
	qb.WriteString("( ")
	values := make([]interface{}, 0, len(wcp.operands))
	for ix, operand := range wcp.operands {
		if ix > 0 {
			qb.WriteString(" or ")
		}
		part, opValues := operand.asWhereClausePart()
		qb.WriteString(part)
		values = append(values, opValues...)
	}
	qb.WriteString(" )")
	return qb.String(), values
}

type wcAnd struct {
	operands []whereClausePart
}

func (wcp *wcAnd) asWhereClausePart() (string, []interface{}) {
	var qb strings.Builder
	qb.WriteString("( ")
	values := make([]interface{}, 0, len(wcp.operands))
	for ix, operand := range wcp.operands {
		if ix > 0 {
			qb.WriteString(" and ")
		}
		part, opValues := operand.asWhereClausePart()
		qb.WriteString(part)
		values = append(values, opValues...)
	}
	qb.WriteString(" )")
	return qb.String(), values
}
