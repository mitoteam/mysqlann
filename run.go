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

func readRow(rows *sql.Rows) (row []interface{}) {
	column_types, _ := rows.ColumnTypes()

	row = make([]interface{}, len(column_types))
	str_data := make([]string, len(column_types))
	pointers := make([]interface{}, len(column_types))

	for i, _ := range column_types {
		pointers[i] = &str_data[i]
	}

	rows.Scan(pointers...)

	for i, ct := range column_types {
		switch t := ct.DatabaseTypeName(); t {
		case "NULL":
			row[i] = nil

		case "VARCHAR":
			row[i] = str_data[i]

		case "BIGINT":
			v, err := strconv.ParseInt(str_data[i], 10, 64)
			if err == nil {
				row[i] = v
			} else {
				row[i] = "[ERROR: " + err.Error() + "]"
			}

		case "DECIMAL":
			v, err := strconv.ParseFloat(str_data[i], 64)
			if err == nil {
				row[i] = v
			} else {
				row[i] = "[ERROR: " + err.Error() + "]"
			}

		default:
			row[i] = "[UNKNOWN DATATYPE: " + t + "]"
		}
	}

	return row
}
