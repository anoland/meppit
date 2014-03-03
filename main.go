package main

import (

    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"    
)
const (
    url = "http://reddit.com/r/kansascity/comments/1w5mel/brunch/.xml"
)



type RSS struct {
    XMLName xml.Name `xml:"rss"`
    Items Items `xml:"channel"`
}
type Items struct {
    XMLName xml.Name `xml:"channel"`
    ItemList []Item `xml:"item"`
}
type Item struct {
    Title string `xml:"Title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
}

func main() {
    client := http.Client{}
    res, err := client.Get(url)

    if err != nil {
        fmt.Println(err)
    }
    body, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close();

    var rss RSS
    xml.Unmarshal(body, &rss)
    if err != nil {
        log.Fatal(err)  
    }

    fmt.Printf("%#v", rss)

    for _, item := range rss.Items.ItemList {
        fmt.Printf("%s\n", item.Description)
    }
    
}
