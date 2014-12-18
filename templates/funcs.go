package templates

import (
  "fmt"
  "html/template"
  "time"
)


func SaneTime(time time.Time) string {
  format := "01/02/2006 15:04"
  return time.Format(format)
}
var funcmap = template.FuncMap{
  "SaneTime": SaneTime,
}
var Templates = template.Must(template.New("").Funcs(funcmap).ParseGlob("templates/*.tpl"))

