package mysqlann

import (
	"database/sql"
	"strconv"
	"strings"
)

type selectQuery struct {
	queryWithWhere

	tables      queryTablesList
	expressions map[string]string
	limit       int
}

func (q *selectQuery) init() {
	q.tables = make(queryTablesList, 0)
}

func Select(table_name string, args ...string) *selectQuery {
	var q selectQuery
	q.init()

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

	table.init()

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

func (q *selectQuery) Expression(expression string, alias string) *selectQuery {
	if q.expressions == nil {
		q.expressions = make(map[string]string, 1)
	}

	q.expressions[alias] = expression

	return q
}

func (q *selectQuery) Limit(limit int) *selectQuery {
	q.limit = limit

	return q
}

func (q *selectQuery) Query() (*sql.Rows, error) {
	return query(q)
}

func (q *selectQuery) QueryRow() (row []interface{}, err error) {
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		row = readRow(rows)
	}

	return row, err
}

func (q *selectQuery) Sql() string {
	var sb strings.Builder
	sb.Grow(1024) //pre-optimization

	//SELECT
	sb.WriteString("SELECT ")

	//fields count
	f_cnt := 0
	for _, table := range q.tables {
		f_cnt += len(table.fields)
	}

	var field_expressions = make(map[string]string, f_cnt+len(q.expressions))

	//FIELDS
	for _, table := range q.tables {
		for field_alias, field_name := range table.fields {
			field_expressions[field_alias] = "`" + table.Alias + "`.`" + field_name + "`"
		}
	}

	//EXPRESSIONS
	if q.expressions != nil {
		for alias, expression := range q.expressions {
			field_expressions[alias] = expression
		}
	}

	first := true
	for alias, expression := range field_expressions {
		if first {
			first = false
		} else {
			sb.WriteString(",\n       ")
		}
		sb.WriteString(expression + " AS `" + alias + "`")
	}

	//FROM
	first = true
	sb.WriteString("\nFROM ")

	if len(q.tables) > 1 {
		sb.WriteString("(")
	}

	for _, table := range q.tables {
		if first {
			first = false
		} else {
			sb.WriteString(", ")
		}

		if table.Name == "DUAL" {
			sb.WriteString(table.Name)
		} else {
			sb.WriteString("`" + table.Name + "` " + table.Alias)
		}
	}

	if len(q.tables) > 1 {
		sb.WriteString(")")
	}

	//WHERE
	q.buildWhere(&sb)

	//LIMIT
	if q.limit > 0 {
		sb.WriteString("\nLIMIT " + strconv.Itoa(q.limit))
	}

	return sb.String()
}
