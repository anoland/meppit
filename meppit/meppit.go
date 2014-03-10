package meppit

import (
    "html/template"
    "net/http"
)

type Page struct {
	Title string
}
func init() {
    http.HandleFunc("/", handler)
}

var templates = template.Must(template.ParseGlob("templates/*.tpl"))


func handler(w http.ResponseWriter, r *http.Request) {
	p := new(Page)
	p.Title = "title" 
	templates.ExecuteTemplate(w, "layout", p)
}
