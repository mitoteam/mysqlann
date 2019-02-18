package mysqlann

import (
	"strings"
)

type deleteQuery struct {
	queryWithWhere

	table_name string
}

func Delete(table_name string) *deleteQuery {
	var q deleteQuery

	q.table_name = table_name

	return &q //for method chaining
}

func (q *deleteQuery) Where(args ...Anything) *deleteQuery {
	q.addWhere(args...)
	return q //method chaining
}

func (q *deleteQuery) Sql() string {
	var sb strings.Builder
	sb.Grow(1024) //pre-optimization

	sb.WriteString("DELETE FROM `")
	sb.WriteString(q.table_name)
	sb.WriteString("`")

	//WHERE
	q.buildWhere(&sb)

	return sb.String()
}
