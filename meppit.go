package main

import (
	"bytes"
    "database/sql"
//	"encoding/json"
    "errors"
	"fmt"
	"html/template"
//	"io/ioutil"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strings"
    "syscall"
	"time"
	_ "github.com/go-sql-driver/mysql"

	"github.com/burntsushi/toml"
    "github.com/mikespook/golib/signal"
)

const (
	url  = "http://reddit.com/r/kansascity/comments/1w5mel/brunch/.json"
	url2 = "http://www.reddit.com/r/kansascity/comments/1ynfb1/who_makes_the_best_reuben_in_town/.json"
)

//var db =
var config Config
var logger = log.New(os.Stdout, "", log.Lshortfile)
type Config struct {
	Version float64 `toml:"version"`
	Host    host    `toml:"host"`
    Database database `toml:"database"`

}
type host struct {
	ListenAddr string `toml:"listenaddr"`
	ListenPort string `toml:"listenport"`
}
type database struct {
    Dbhost string `toml:"dbhost"`
    Dbuser string `toml:"dbuser"`
    Dbpass string `toml:"dbpass"`
    Dbname string `toml:"dbname"`
}

func main() {
    config := Config{}
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal("Problem with config file", err)
	}

	logger.Println("Running version:", config.Version)

    database := config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", database.Dbuser, database.Dbpass, database.Dbhost, database.Dbname)
	db, _ = sql.Open("mysql", dsn)
	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
    if err := setup(); err != nil {
        logger.Fatal(err)
    }


	http.HandleFunc("/", indexHandler)
    http.HandleFunc("/css/", staticHandler)
    http.HandleFunc("/js/", staticHandler)
    http.HandleFunc("/fetch", fetchHandler)
	listen := fmt.Sprintf("%s:%s", config.Host.ListenAddr, config.Host.ListenPort)
	logger.Println("Starting server on: ", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		logger.Fatal("Problem starting server", err)
		return
	}

	signal.Bind(syscall.SIGINT, func() uint { return signal.BreakExit})
    s := signal.Wait()
	logger.Printf("Exit by signal: %s\n", s)

}

func staticHandler(w http.ResponseWriter, r *http.Request) {
   http.ServeFile(w, r, filepath.Join("static", r.URL.Path)) 

}
func indexHandler(w http.ResponseWriter, r *http.Request) {
    
    
    p := struct{
        Title string
        Content string
    }{
        "index page",
        "index page contnet",
    }
    RenderTemplate(w, "index", p)
}
func fetchHandler(w http.ResponseWriter, r *http.Request) {
    p := struct {
        Title string
        Content string

    }{
       "fetch page",
       "fetch page",
    }
    RenderTemplate(w, "fetch", p)
}

func setup () error{
    if err := setupDB(); err != nil {
        return err
    }
    if err := setupTemplates(); err != nil {
        return err
    }
    return nil

}
var db *sql.DB
func setupDB () error{
    logger.Println("running setup")
    logger.Println("setting up places table")
    places_sql := `
    CREATE TABLE IF NOT EXISTS places(
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(20) NOT NULL,
        description VARCHAR(255),
        lat FLOAT,
        lon FLOAT
        );`
    if _, err := db.Exec(places_sql); err != nil {
        return err
    }
    logger.Println("setting up users table")
    users_sql := `
    CREATE TABLE IF NOT EXISTS users(
        id INT PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(20),
        token VARCHAR(200)
    );`
    if _, err := db.Exec(users_sql); err != nil {
        return err
    }
    logger.Println("setup complete")
    return nil
}
var templates map[string]*template.Template
func setupTemplates() error{
    if templates == nil {
        templates = make(map[string]*template.Template)
    }


    pages, err := filepath.Glob("templates/*.tpl") 
    if err != nil {
       return err
    }

    for _, page := range pages {
        name := strings.TrimSuffix(filepath.Base(page), filepath.Ext(page))
        templates[name] = template.Must(template.ParseFiles("layout.tpl", page))
    }
    return nil
}
func RenderTemplate(w http.ResponseWriter, name string, p interface{}) error {

    if p == nil {
        err := errors.New("could not find page")
         http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }
    
    tmpl, ok := templates[name]
    if !ok {
        logger.Println(fmt.Errorf("The template %s does not exist.", name))
    }

    var buff bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buff, "layout", p); err != nil {
        logger.Println("Error executing template: ", name, err)
    }
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    buff.WriteTo(w)

    return nil
}
func SaneTime(time time.Time) string {
	format := "01/02/2006 15:04"
	return time.Format(format)
}

var funcmap = template.FuncMap{
	"SaneTime": SaneTime,
}

