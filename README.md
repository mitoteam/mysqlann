# MiTo Team MysqlAnn
MySQL query builder for Golang

*WARNING*: Not ready for production yet

* Building SQL queries
* Executing SQL queries

```go
package main

import (
  "fmt"
  "github.com/mitoteam/mysqlann"
  "database/sql"
)

func main(){
	//creating query
	q := mysqlann.Select("mt_user", "u", "Role", "UserName").
		Where("ID", 67).
		Limit(1)
	fmt.Println("SQL:")
	fmt.Println(q.Sql())

	//initialize connection pool
	db, _ := sql.Open("mysql", "user:password@tcp(webdev4.test)/test") //use your own DSN
	mysqlann.SetDB(db)

	//run query
	row_values, _ := q.QueryRowMap()
	fmt.Println("Result: ", row_values)
}
```

Will produce

```
SQL:
SELECT `u`.`Role` AS `Role`,
       `u`.`UserName` AS `UserName`
FROM `mt_user` u
WHERE (ID = 67)
LIMIT 1

Result:  map[Role:admin UserName:TestUser]
```
