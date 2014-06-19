package templates

import (
  "fmt"
  "html/template"
  "time"
 // "strconv"
  "labix.org/v2/mgo/bson"
)


func StringId(id bson.ObjectId) string {
  i := id.Hex()
  return i 
}
func SaneTime(time time.Time) string {
  format := "01/02/2006 15:04"
  return time.Format(format)
}
func PrettyPrint(phone string) template.HTML {
  return template.HTML(fmt.Sprintf("<div class='phone'>%s</div>", phone))
}
var funcmap = template.FuncMap{
  "StringId": StringId,
  "SaneTime": SaneTime,
  "PrettyPrint": PrettyPrint,
}
var Templates = template.Must(template.New("").Funcs(funcmap).Delims("[[", "]]").ParseGlob("templates/*.tpl"))

