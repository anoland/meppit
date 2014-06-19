package main

import (
  "net/http"
    
  "anoland.rhcloud.com/meppit/domain"
  "anoland.rhcloud.com/meppit/templates"
)

func init() {
    http.HandleFunc("/user/", userHandler)
}


func userHandler(w http.ResponseWriter, r *http.Request) {
  sess := domain.GetSession()
  defer sess.Close()
  
  offset := 10
  
  coll := sess.DB(databaseName).C("people")
  count, _ := coll.Count()
  
  results := []*domain.Person{}
  rs := coll.Find(nil).Limit(10).Skip(offset)
  err := rs.All(&results)
  
  if err != nil {
    panic(err)
  }
  
  content := struct {
    Count int
    Results []*domain.Person
  }{
    count,
    results,
  }
  
  templates.Templates.ExecuteTemplate(w, "mgo", content)
  
}
