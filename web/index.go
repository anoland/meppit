package web

import (
    "net/http"

)


func init() {
    http.HandleFunc("/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}
