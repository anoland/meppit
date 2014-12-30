package web

import (
    "net/http"

)


func init() {
    http.HandleFunc("/submit", submitHandler)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello")
}
