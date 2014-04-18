package main

import (
	"reflect"
	"strconv"
	"text/template"
)

var tplFuncs = map[string]interface{}{
	"length": func(param interface{}) int {
		return reflect.ValueOf(param).Len()
	},
	"quote": strconv.Quote,
}

const structTemplateStr = `
func (self *{{.Name}}) BindFields(names []string, values []interface{}) {
	for i, name := range names {
		switch name {
{{range.Fields}}
		case {{quote .SqlName}}:
			values[i] = &self.{{.Name}}
{{end}}
		}
	}
}

func (self *{{.Name}}) EnumerateFields() (names []string, values []interface{}) {
	names = make([]string, 0, {{length .Fields}})
	values = make([]interface{}, 0, {{length .Fields}})
{{range .Fields}}
	names = append(names, {{quote .SqlName}})
	values = append(values, &self.{{.Name}})
{{end}}
	return
}
`

var structTemplate = template.Must(template.New("").Funcs(tplFuncs).Parse(structTemplateStr))
