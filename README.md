# MiTo Team MysqlAnn
MySQL query builder for Golang

```go
	var q = mysqlann.Select("mt_system_users", "u", "ID", "UserName")
	fmt.Println(q.Sql())
```

Will produce

```mysql
    SELECT `u`.`UserName` AS `UserName`,
           `u`.`ID` AS `ID`
    FROM (`mt_system_users` u)
```