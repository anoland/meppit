package meppit

import (

    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"    



)
const (
    url = "http://reddit.com/r/kansascity/comments/1w5mel/brunch/.xml"
    url2 = "http://www.reddit.com/r/kansascity/comments/1ynfb1/who_makes_the_best_reuben_in_town/.xml"
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

func init() {
    http.HandleFunc("/fetch", fetchHandler)
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
// re-insert http client to grab url
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
	items := make(map[int]Item)
    for i, item := range rss.Items.ItemList {
	[i]items := item
	fmt.Fprintf(w, "%s \r\n", item.Description)
    }
    
	p := new(Page)
	p.Title = "title" 
	templates.ExecuteTemplate(w, "fetch", p)
}
