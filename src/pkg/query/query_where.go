package query

import (
	"fmt"
	"strings"
)

type WhereClausePart interface {
	AsWhereClausePart() (string, []interface{})
}

func (qc *QualifiedColumn) AsWhereClausePart() (string, []interface{}) {
	return fmt.Sprintf("%s.%s = $$", qc.TableName, qc.ColumnName), []interface{}{qc.Value}
}

type WcOr struct {
	Operands []WhereClausePart
}

func (wcp *WcOr) AsWhereClausePart() (string, []interface{}) {
	var qb strings.Builder
	qb.WriteString("( ")
	values := make([]interface{}, 0, len(wcp.Operands))
	for ix, operand := range wcp.Operands {
		if ix > 0 {
			qb.WriteString(" or ")
		}
		part, opValues := operand.AsWhereClausePart()
		qb.WriteString(part)
		values = append(values, opValues...)
	}
	qb.WriteString(" )")
	return qb.String(), values
}

type WcAnd struct {
	Operands []WhereClausePart
}

func (wcp *WcAnd) AsWhereClausePart() (string, []interface{}) {
	var qb strings.Builder
	qb.WriteString("( ")
	values := make([]interface{}, 0, len(wcp.Operands))
	for ix, operand := range wcp.Operands {
		if ix > 0 {
			qb.WriteString(" and ")
		}
		part, opValues := operand.AsWhereClausePart()
		qb.WriteString(part)
		values = append(values, opValues...)
	}
	qb.WriteString(" )")
	return qb.String(), values
}
