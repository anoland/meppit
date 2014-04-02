package meppit

import (
    "html/template"
    "net/http"
)

type Title string

func init() {
    http.HandleFunc("/", handler)
}

var templates = template.Must(template.ParseGlob("meppit/templates/*.tpl"))

func handler(w http.ResponseWriter, r *http.Request) {
	index := new(Page)
	index.SetTitle("title") 
	templates.ExecuteTemplate(w, "index", index)
}
