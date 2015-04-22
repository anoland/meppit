package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/jzelinskie/geddit"
	"gopkg.in/unrolled/render.v1"
)

const (
	url  = "http://www.reddit.com/r/kansascity/comments/1w5mel/brunch/.json"
	url2 = "http://www.reddit.com/r/kansascity/comments/1ynfb1/who_makes_the_best_reuben_in_town/.json"
)

var conn *sql.DB
var config Config
var r *render.Render
var logger = log.New(os.Stdout, "", log.Lshortfile)

type Config struct {
	Version  float64  `json:"version"`
	Host     host     `json:"host"`
	Database database `json:"database"`
}
type host struct {
	ListenAddr string `json:"listenaddr"`
	ListenPort string `json:"listenport"`
}
type database struct {
	Dbhost string `json:"dbhost"`
	Dbuser string `json:"dbuser"`
	Dbpass string `json:"dbpass"`
	Dbname string `json:"dbname"`
}

func main() {
	cf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := Config{}
	if err := json.Unmarshal(cf, &config); err != nil {
		logger.Fatal("Problem with config file", err)
	}

	logger.Println("Running version:", config.Version)

	dbconfig := config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbconfig.Dbuser, dbconfig.Dbpass, dbconfig.Dbhost, dbconfig.Dbname)
	conn, _ = sql.Open("mysql", dsn)
	if err := conn.Ping(); err != nil {
		logger.Fatal(err)
	}
	var debug bool
	app_env := os.Getenv("APP_ENVIRONMENT")
	logger.Println("Running in environment: ", app_env)
	if app_env == "development" {
		debug = true
		logger.Println("debugging is on")
	} else {
		debug = false
		logger.Println("debugging is off")
	}
	r = render.New(render.Options{
		Directory:     "templates",
		Layout:        "layout",
		IsDevelopment: debug,
	})

	go func() {
		fetchJob()
		ticker := time.Tick(1 * time.Hour)
		for range ticker {
			fetchJob()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/admin/fetch/", fetchHandler)

	listen := fmt.Sprintf("%s:%s", config.Host.ListenAddr, config.Host.ListenPort)
	logger.Println("Starting server on: ", listen)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		logger.Fatal("Problem starting server", err)
		return
	}

}
func indexHandler(w http.ResponseWriter, req *http.Request) {
	p := struct {
		Title string
		Body  string
	}{
		"Index page",
		"This is the body",
	}
	r.HTML(w, http.StatusOK, "index", p)
}

func fetchHandler(w http.ResponseWriter, req *http.Request) {
	var (
		title     string
		reddit_id string
		url       string
	)
	type Submission struct {
		Title    string
		RedditID string
		URL      string
	}

	var submissions []*Submission
	rs, err := conn.Query("select title, reddit_id, url from submissions")
	if err != nil {
		logger.Println("fetching rows failed")
	}
	defer rs.Close()
	for rs.Next() {
		err := rs.Scan(&title, &reddit_id, &url)
		if err != nil {
			logger.Println(err)
		}
		s := Submission{
			title,
			reddit_id,
			url,
		}
		submissions = append(submissions, &s)
	}
	err = rs.Err()
	if err != nil {
		logger.Println(err)
	}
	p := struct {
		Title       string
		Submissions []*Submission
	}{
		"fetch page",
		submissions,
	}
	r.HTML(w, http.StatusOK, "admin/fetch", p)
}

func fetchJob() {
	logger.Println("starting fetch job")

	fetch_sql, err := conn.Prepare("insert into submissions (url, reddit_id, title) values ( ?, ?, ?)")
	if err != nil {
		logger.Println(err)
	}

	rs, err := conn.Query("select reddit_id from submissions")
	if err != nil {
		logger.Println(err)
	}
	defer rs.Close()
	var haves []string
	for rs.Next() {
		var have string
		err := rs.Scan(&have)
		if err != nil {
			logger.Println(err)
		}
		haves = append(haves, have)
	}
	err = rs.Err()
	if err != nil {
		logger.Fatal(err)
	}
	subreddit := "golang"
	// todo: add toggle meppitdev
	ua := fmt.Sprintf("web:Meppit: /r/meppit -  maps for reddit -:v %s (by /u/anoland)", config.Version)
	ggg := geddit.NewSession(ua)
	submissions, _ := ggg.SubredditSubmissions(subreddit)
	count := 0
	for _, sub := range submissions {
		_, ok := idExists(haves, sub.ID)
		if !ok {
			_, err := fetch_sql.Exec(sub.URL, sub.ID, sub.Title)
			count = count + 1
			if err != nil {
				logger.Println(err)
			}
		}

	}
	// todo: keep track of last sucessful fetch

	logger.Println("finished fetch job.")
	logger.Printf("processed %d submissions", count)
}

func idExists2(haystack []string, needle string) (int, bool) {
	sort.Strings(haystack)
	l := len(haystack)
	for i := 0; i < l; i++ {
		value := haystack[i]
		if value < needle {
			continue
		}
		return i, value == needle
	}
	return l, false
}

// reconfiguration of above using binary (builtin) search instead
func idExists(haystack []string, needle string) (int, bool) {
	sort.Strings(haystack)
	l := len(haystack)
	i := sort.SearchStrings(haystack, needle)
	if i < l && haystack[i] == needle {
		return i, true
	}
	return l, false

}
