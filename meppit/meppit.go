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

func handler(w http.ResponseWriter, r *http.Request) {
	p := Page{"title"}
	t := template.Must(template.ParseGlob("meppit/templates/*"))
	t.ExecuteTemplate(w, "layout", p)
}
