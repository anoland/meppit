package main

import (
    "fmt"
	"log"
    //"net"
	"net/http"

    "github.com/burntsushi/toml"
)
	
type Config struct {
    Version float64 `toml:"version"`
    Host host `toml:"host"`
    
}
type host struct {
    ListenAddr string `toml:"listenaddr"`
    ListenPort string `toml:"listenport"`

    }
func init() {
    http.HandleFunc("/", indexHandler)
    }

func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("helo"))

    }
func main() {
    var config Config
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatal("Problem with config file", err)
    }

    listen := fmt.Sprintf("%s:%s", config.Host.ListenAddr, config.Host.ListenPort)

    if err := http.ListenAndServe(listen, nil); err != nil {
        log.Fatal("Problem starting server", err)
        return
    } 

    log.Println("Started server on ", listen)
    log.Println("Running version:", config.Version)
    
}
