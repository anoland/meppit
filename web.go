package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
  
  "anoland.rhcloud.com/meppit/domain"
  "anoland.rhcloud.com/meppit/templates"
  
)

const layout = "10/02/14 3:04pm"
const	databaseName = "test"



func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/mgo", mgoHandler)
  http.HandleFunc("/mgo/edit", mgoAddHandler)
}
func main() {
	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	log.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}



func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello, World from %s", runtime.Version())
}

func mgoHandler(w http.ResponseWriter, req *http.Request) {
	sess := domain.GetSession()
	defer sess.Close()

  coll := sess.DB(databaseName).C("people")
  
  offset := 10
	results := []*domain.Person{}
  iter := coll.Find(nil).Limit(10).Skip(offset)
  count, _ := coll.Count()
  err := iter.All(&results)
	if err != nil {
		panic(err)
	}

  println(count);
	 templates.Templates.ExecuteTemplate(w, "mgo", results)

}


func mgoAddHandler(w http.ResponseWriter, req *http.Request) {
  	session := domain.GetSession()
	defer session.Close()

	c := session.DB(databaseName).C("people")
  p := domain.NewPerson("Betty", "888 098 876")
  err := c.Insert(&p)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
  
}
