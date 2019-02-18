package mysqlann

import "strings"

type baseQuery struct {
	initialized bool
	conditions  *queryCondition
}

func (q *baseQuery) init() {
	q.conditions = And()
}

func (q *baseQuery) checkInit() {
	if !q.initialized {
		q.init()
		q.initialized = true
	}
}

func (q *baseQuery) addWhere(args ...Anything) {
	q.checkInit()

	q.conditions.Add(args...)
}

func (q *baseQuery) buildWhere(sb *strings.Builder) {
	//not initialized
	if q.conditions == nil {
		return
	}

	//no conditions added
	if q.conditions.IsEmpty() {
		return
	}

	sb.WriteString("\nWHERE ")
	sb.WriteString(q.conditions.Sql(true))
}

//https://gist.github.com/siddontang/8875771
func Escape(value string) string {
	dest := make([]byte, 0, 2*len(value))
	var escape byte
	for i := 0; i < len(value); i++ {
		c := value[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}