package pages

import (
    "bytes"
    "fmt"
    "log"
	"html/template"
	"net/http"
	"time"
)

func SaneTime(time time.Time) string {
	format := "01/02/2006 15:04"
	return time.Format(format)
}

var funcmap = template.FuncMap{
	"SaneTime": SaneTime,
}

// RenderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, p *Page) error {
    tmpl := templates.Must(templates.New("layout").ParseFiles("pages/layout.tpl")

    var buff bytes.Buffer

	err := tmpl.ExecuteTemplate(&buff, "base", p)
    if err != nil {
        log.Println("error executing template ", err)
    }

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    buff.WriteTo(w)
	return nil
}

