package mysqlann

import "strings"

type queryWithFieldsValues struct {
	initialized  bool
	field_values queryFieldsValues
}

func (q *queryWithFieldsValues) init() {
	q.field_values = make(queryFieldsValues)
}

func (q *queryWithFieldsValues) checkInit() {
	if !q.initialized {
		q.init()
		q.initialized = true
	}
}

func (q *queryWithFieldsValues) setFieldValue(field_name string, value Anything) {
	q.checkInit()

	q.field_values[field_name] = value
}

func (q *queryWithFieldsValues) buildFieldsValues(sb *strings.Builder, insert bool) {
	//not initialized
	if !q.initialized {
		return
	}

	//no field_values added
	if len(q.field_values) == 0 {
		return
	}

	if insert {
		sb.WriteString("(")

		var first= true
		for field_name, _ := range q.field_values {
			if first {
				first = false
			} else {
				sb.WriteString(", ")
			}

			sb.WriteString(field_name)
		}

		sb.WriteString(")\nVALUES (")

		first = true
		for _, value := range q.field_values {
			if first {
				first = false
			} else {
				sb.WriteString(", ")
			}

			sb.WriteString(AnythingToSql(value))
		}

		sb.WriteString(")")
	}
}
