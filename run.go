package mysqlann

import (
	"database/sql"
	"strconv"
)

var db *sql.DB

func SetDB(new_db *sql.DB) {
	db = new_db
}

func exec(q Query) (sql.Result, error) {
	if db == nil {
		panic("connection not initialized")
	}

	/*
	stmt, _ := db.Prepare(q.Sql())
	defer stmt.Close()
	return stmt.Exec()
	*/

	return db.Exec(q.Sql())
}

func query(q Query) (*sql.Rows, error) {
	if db == nil {
		panic("connection not initialized")
	}

	/*
	stmt, _ := db.Prepare(q.Sql())
	defer stmt.Close()
	return stmt.Query()
	*/

	return db.Query(q.Sql())
}

func readRow(rows *sql.Rows) (row []Anything) {
	column_types, _ := rows.ColumnTypes()

	row = make([]Anything, len(column_types))
	str_data := make([]string, len(column_types))
	pointers := make([]interface{}, len(column_types))

	for i, _ := range column_types {
		pointers[i] = &str_data[i]
	}

	rows.Scan(pointers...)

	for i, columnType := range column_types {
		row[i] = strToValue(str_data[i], columnType)
	}

	return row
}

func readSingleColumnByIndex(rows *sql.Rows, columnIndex int) (value Anything) {
	column_types, _ := rows.ColumnTypes()
	var column_type *sql.ColumnType

	if len(column_types) <= columnIndex {
		return nil //no such column
	} else {
		column_type = column_types[columnIndex]
	}

	var str_data string
	rows.Scan(&str_data)

	value = strToValue(str_data, column_type)

	return value
}

func readRowMap(rows *sql.Rows) (row map[string]interface{}) {
	column_types, _ := rows.ColumnTypes()

	row = make(map[string]interface{}, len(column_types))
	str_data := make([]string, len(column_types))
	pointers := make([]interface{}, len(column_types))

	for i, _ := range column_types {
		pointers[i] = &str_data[i]
	}

	rows.Scan(pointers...)

	for i, ct := range column_types {
		row[ct.Name()] = strToValue(str_data[i], ct)
	}

	return row
}

func strToValue(str string, ct *sql.ColumnType) (r interface{}) {
	switch t := ct.DatabaseTypeName(); t {
	case "NULL":
		r = nil

	case "VARCHAR":
		r = str

	case "INT":
		v, err := strconv.Atoi(str)
		if err == nil {
			r = v
		} else {
			r = "[ERROR: " + err.Error() + "]"
		}

	case "BIGINT":
		v, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			r = v
		} else {
			r = "[ERROR: " + err.Error() + "]"
		}

	case "DECIMAL":
		v, err := strconv.ParseFloat(str, 64)
		if err == nil {
			r = v
		} else {
			r = "[ERROR: " + err.Error() + "]"
		}

	default:
		r = "[UNKNOWN DATATYPE: " + t + "]"
	}

	return r
}
