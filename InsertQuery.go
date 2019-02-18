package mysqlann

import (
	"strings"
)

type insertQuery struct {
	queryWithFieldsValues

	table_name string
}

func Insert(table_name string) *insertQuery {
	var q insertQuery

	q.table_name = table_name

	return &q //for method chaining
}

func (q *insertQuery) Set(field_name string, value Anything) *insertQuery{
	q.setFieldValue(field_name, value)
	return q //method chaining
}

func (q *insertQuery) Sql() string {
	var sb strings.Builder
	sb.Grow(1024) //pre-optimization

	sb.WriteString("INSERT INTO `")
	sb.WriteString(q.table_name)
	sb.WriteString("`")

	//FIELDS AND VALUES
	q.buildFieldsValues(&sb, true)

	return sb.String()
}
