package mysqlann

import "strings"

type ConditionGlueOperator int
type queryConditionList []*queryCondition

//glue types
const (
	OR  = 1
	AND = 2
)

type queryCondition struct {
	operator          ConditionGlueOperator
	string_conditions []string
	sub_conditions    queryConditionList
}

func And() *queryCondition {
	qc := &queryCondition{
		operator: AND,
	}

	qc.init()

	return qc
}

func Or() *queryCondition {
	qc := &queryCondition{
		operator: OR,
	}

	qc.init()

	return qc
}

func (qc *queryCondition) init() {
	qc.string_conditions = make([]string, 0, 1) //assuming there will be at least one simple string condition
	qc.sub_conditions = make(queryConditionList, 0)
}

func (qc *queryCondition) IsEmpty() bool {
	return len(qc.string_conditions) == 0 && len(qc.sub_conditions) == 0
}

func (qc *queryCondition) AddString(condition string) *queryCondition {
	qc.string_conditions = append(qc.string_conditions, condition)

	return qc //method chaining
}

func (qc *queryCondition) AddEquality(field string, value Anything) *queryCondition {
	if value == nil {
		qc.AddString(field + " IS NULL")
	} else {
		qc.AddWithOperator(field, "=", value)
	}

	return qc //method chaining
}

func (qc *queryCondition) AddWithOperator(field string, operator string, value Anything) *queryCondition {
	var str = field + " " + operator + " " + AnythingToSql(value)

	qc.AddString(str)

	return qc //method chaining
}

func (qc *queryCondition) AddCondition(condition *queryCondition) *queryCondition {
	qc.sub_conditions = append(qc.sub_conditions, condition)

	return qc //method chaining
}

func (qc *queryCondition) Add(args ...Anything) *queryCondition {
	if len(args) == 0 || len(args) > 3 {
		panic("There should be 1, 2 or 3 args here")
	}

	if len(args) == 1 {
		condition, isCondition := args[0].(*queryCondition)

		if isCondition {
			qc.AddCondition(condition)
		} else {
			qc.AddString(AnythingToString(args[0]))
		}
	} else if len(args) == 2 {
		qc.AddEquality(args[0].(string), args[1])
	} else if len(args) == 3 {
		qc.AddWithOperator(args[0].(string), args[1].(string), args[2])
	}

	return qc //method chaining
}

func (qc *queryCondition) Sql(multiline bool) string {
	//preapare conditions as strings
	var cooked_string_conditions = make([]string, len(qc.string_conditions) + len(qc.sub_conditions))

	i := 0
	for _, string_condition:= range qc.string_conditions {
		cooked_string_conditions[i] = "(" + string_condition + ")"
		i++
	}

	i = len(qc.string_conditions)
	for _, condition:= range qc.sub_conditions {
		cooked_string_conditions[i] = "(" + condition.Sql(false) + ")"
		i++
	}

	//join into one string
	var sb strings.Builder

	first := true

	for _, string_condition := range cooked_string_conditions {
		if first {
			first = false
		} else {
			if multiline{
				sb.WriteString("\n  ")
			} else{
				sb.WriteString(" ")
			}
		}

		sb.WriteString(qc.OperatorSql())
		sb.WriteString(" ")
		sb.WriteString(string_condition)
	}

	return sb.String()
}

func (qc *queryCondition) OperatorSql() string {
	switch qc.operator {
	case OR:
		return "OR"
	case AND:
		return "AND"
	default:
		return "[unknown operator]"
	}
}
