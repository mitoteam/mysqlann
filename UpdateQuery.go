package mysqlann

import (
	"strings"
)

type updateQuery struct {
	queryWithFieldsValues
	queryWithWhere

	table_name string
}

func Update(table_name string) *updateQuery {
	var q updateQuery

	q.table_name = table_name

	return &q //for method chaining
}

func (q *updateQuery) Set(field_name string, value Anything) *updateQuery{
	q.setFieldValue(field_name, value)
	return q //method chaining
}

func (q *updateQuery) Where(args ...Anything) *updateQuery {
	q.addWhere(args...)
	return q //method chaining
}

func (q *updateQuery) Sql() string {
	var sb strings.Builder
	sb.Grow(1024) //pre-optimization

	sb.WriteString("UPDATE `")
	sb.WriteString(q.table_name)
	sb.WriteString("`\nSET ")

	//FIELDS AND VALUES
	q.buildFieldsValues(&sb, false)

	//WHERE
	q.buildWhere(&sb)

	return sb.String()
}
