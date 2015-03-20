package main

import (
    "fmt"
	"log"
	"net/http"

    "github.com/burntsushi/toml"

    _ "github.com/anoland/meppit/web"
)
	
var config Config
type Config struct {
    Version float64 `toml:"version"`
    Host host `toml:"host"`
    
}
type host struct {
    ListenAddr string `toml:"listenaddr"`
    ListenPort string `toml:"listenport"`

}

func main() {
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatal("Problem with config file", err)
    }

    listen := fmt.Sprintf("%s:%s", config.Host.ListenAddr, config.Host.ListenPort)

    if err := http.ListenAndServe(listen, nil); err != nil {
        log.Fatal("Problem starting server", err)
        return
    } 

    fmt.Println("Started server on ", listen)
    fmt.Println("Running version:", config.Version)
    
}
