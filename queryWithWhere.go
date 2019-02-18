package mysqlann

import "strings"

type queryWithWhere struct {
	initialized     bool
	where_condition *queryCondition
}

func (q *queryWithWhere) init() {
	q.where_condition = And()
}

func (q *queryWithWhere) checkInit() {
	if !q.initialized {
		q.init()
		q.initialized = true
	}
}

func (q *queryWithWhere) addWhere(args ...Anything) {
	q.checkInit()

	q.where_condition.Add(args...)
}

func (q *queryWithWhere) buildWhere(sb *strings.Builder) {
	//not initialized
	if !q.initialized {
		return
	}

	//no conditions added
	if q.where_condition.IsEmpty() {
		return
	}

	sb.WriteString("\nWHERE ")
	sb.WriteString(q.where_condition.Sql(true))
}
