// +build !appengine

package main

import (
	"log"
	"net/http"
)
	

func main() {
	err := http.ListenAndServe("10.5.1.54:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
