package template

import "html/template"

var All = template.Must(template.ParseFiles(
	"template/test.html",
	"template/list-tutorial.html",
	"template/edit-tutorial.html",
	"template/add-tutorial.html"))
