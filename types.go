package mysqlann

import (
	"fmt"
	"strconv"
)

//region Anything
type Anything interface{}

func AnythingToSql(value Anything) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	default: //strings are also here
		return "'" + Escape(AnythingToString(value)) + "'"
	}
}

func AnythingToString(value Anything) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}
//endregion

//region Query interface
type Query interface {
	Sql() string
}
//endregion

type queryTablesList []*queryTable

type queryFieldsValues map[string]Anything
