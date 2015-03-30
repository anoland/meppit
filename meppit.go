package main

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "time"

    "github.com/burntsushi/toml"
    _ "github.com/anoland/meppit/web"
)

const (
    url = "http://reddit.com/r/kansascity/comments/1w5mel/brunch/.xml"
    url2 = "http://www.reddit.com/r/kansascity/comments/1ynfb1/who_makes_the_best_reuben_in_town/.xml"
)

//var db = 
var tpl = template.Must(template.New("").Funcs(funcmap).ParseGlob("templates/*.tpl"))
var config Config
type Config struct {
    Version float64 `toml:"version"`
    Host host `toml:"host"`
}
type host struct {
    ListenAddr string `toml:"listenaddr"`
    ListenPort string `toml:"listenport"`

}


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
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/fetch", fetchHandler)
    http.HandleFunc("/submit", submitHandler)
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



type Page struct {
    Title string
    Content string 
}


func NewPage() *Page {
    p := &Page{}
    return p

}

func (p *Page) SetTitle(title string) (error) {
    p.Title = title
    return nil
}

func (p *Page) SetContent(content string) (error) {
    p.Content = content
    return nil
}
func SaneTime(time time.Time) string {
  format := "01/02/2006 15:04"
  return time.Format(format)
}
var funcmap = template.FuncMap{
  "SaneTime": SaneTime,
}

// RenderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, p *Page) error {
    t := tpl.Must(templates.New("layout").ParseFiles("templates/layout.tpl")

    var buff bytes.Buffer

	err := t.ExecuteTemplate(&buff, "base", p)
    if err != nil {
        log.Println("error executing template ", err)
    }

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    buff.WriteTo(w)
	return nil
}





func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
    res, err := http.Get(url)
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
    for _, item := range rss.Items.ItemList {
        fmt.Fprintf(w, "%s \r\n<br />", item.Description)
    }
    
}


func submitHandler(w http.ResponseWriter, r *http.Request) {
    submit := NewPage()
    submit.SetTitle("Submit") 
    submit.SetContent("this is content")
}
