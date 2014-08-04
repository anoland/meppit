package meppit

import (

    "net/http"    

  //  "appengine"

)

func init() {
    http.HandleFunc("/submit", submitHandler)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    submit := NewPage()
    submit.SetTitle("Submit") 
    submit.SetContent("this is content")
    templates.ExecuteTemplate(w, "submit", submit)
//    ctx := appengine.NewContext(r)
}
