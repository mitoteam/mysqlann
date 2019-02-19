package mysqlann

type tableFieldsList map[string]string

type queryTable struct {
	Name string
	Alias string

	all bool
	fields tableFieldsList
}

func (q *queryTable) AddField(field_name string, field_alias string){
	if q.fields == nil {
		q.fields = make(tableFieldsList, 1)
	}

	if len(field_alias) == 0 {
		field_alias = field_name
	}

	q.fields[field_alias] = field_name
}
