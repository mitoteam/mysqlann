package mysqlann

type tableFieldsList map[string]string

type queryTable struct {
	Name string
	Alias string

	fields tableFieldsList
}

func (q *queryTable) init(){
	q.fields = make(tableFieldsList)
}

func (q *queryTable) AddField(field_name string, field_alias string){
	if len(field_alias) == 0 {
		field_alias = field_name
	}

	q.fields[field_alias] = field_name
}
