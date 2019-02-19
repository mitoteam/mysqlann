package mysqlann

import "database/sql"

var db *sql.DB

func SetDB(new_db *sql.DB) {
	db = new_db
}

func exec(q Query) (sql.Result, error){
	if(db == nil){
		panic("connection not initialized")
	}

	return db.Exec(q.Sql())
}

func query(q Query) (sql.Rows, error){
	if(db == nil){
		panic("connection not initialized")
	}

	return db.Query(q.Sql())
}
