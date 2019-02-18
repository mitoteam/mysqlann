package mysqlann

import (
	"strconv"
	"strings"
)

type queryTablesList []*queryTable

type selectQuery struct {
	baseQuery

	tables queryTablesList
}

func (q *selectQuery) init() {
	q.tables = make(queryTablesList, 0)
}

func Select(table_name string, args ...string) *selectQuery {
	var q selectQuery
	q.init();

	alias := ""

	if len(args) > 0 {
		alias = args[0]
	}

	if len(args) > 1 {
		var fields = args[1:]
		q.AddTable(table_name, alias, fields...)
	} else {
		q.AddTable(table_name, alias)
	}

	return &q //for method chaining
}

func (q *selectQuery) AddTable(table_name string, alias string, fields ...string) *selectQuery {
	if len(alias) == 0 {
		alias = "t" + strconv.Itoa((len(q.tables) + 1))
	}

	table := &queryTable{
		Name:  table_name,
		Alias: alias,
	}

	table.init();

	for _, field_name := range fields {
		table.AddField(field_name, "")
	}

	q.tables = append(q.tables, table)

	return q //for method chaining
}

func (q *selectQuery) Where(args ...Anything) *selectQuery {
	q.addWhere(args...)

	return q
}

func (q *selectQuery) Sql() string {
	var sb strings.Builder
	sb.Grow(1024) //pre-optimization

	sb.WriteString("SELECT ")

	//FIELDS
	first := true
	for _, table := range q.tables {
		for field_alias, field_name := range table.fields {
			if(first){
				first = false
			} else {
				sb.WriteString(",\n       ")
			}
			sb.WriteString("`" + table.Alias + "`.`" + field_name + "` AS `" + field_alias + "`")
		}
	}

	//FROM
	sb.WriteString("\nFROM (")
	for _, table := range q.tables {
		sb.WriteString("`" + table.Name + "` " + table.Alias)
	}
	sb.WriteString(")")

	//WHERE
	q.buildWhere(&sb)

	return sb.String()
}
