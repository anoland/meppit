package main

import (
	//_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
    
	"github.com/anoland/geddit"
	"github.com/codegangsta/negroni"
   "github.com/julienschmidt/httprouter" 
    "github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
)

const (
	url  = "http://www.reddit.com/r/kansascity/comments/1w5mel/brunch/.json"
	url2 = "http://www.reddit.com/r/kansascity/comments/1ynfb1/who_makes_the_best_reuben_in_town/.json"
)

var conn *sqlx.DB
var config *Config
var r *render.Render
var logger = log.New(os.Stdout, "", log.Lshortfile)
var debug bool 

type Config struct {
	Version  float64  `json:"version"`
    Reddituser string `json:"reddituser"`
    Redditpass string `json:"redditpass"`
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
	conn, _ = sqlx.Open("mysql", dsn)
	if err := conn.Ping(); err != nil {
		logger.Fatal(err)
	}
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

	go fetchJob(config)

	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/admin/fetch/", fetchHandler)
	router.GET("/admin/fetch/detail/:id", fetchDetailHandler)

	listen := fmt.Sprintf("%s:%s", config.Host.ListenAddr, config.Host.ListenPort)
	logger.Println("Starting server on: ", listen)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		logger.Fatal("Problem starting server", err)
		return
	}

}

func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p := struct {
		Title string
		Body  string
	}{
		"Index page",
		"This is the body",
	}
	r.HTML(w, http.StatusOK, "index", p)
}


func fetchHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

    type Submission struct {
        Title string
        Permalink string
        RedditID string  `db:"reddit_id"`
        URL string
    }
    submissions := []Submission{}
	err := conn.Select(&submissions, "select title, permalink, reddit_id, url from submissions")
	if err != nil {
		logger.Println(err)
	}
    p := struct {
        Title string
        Submissions []Submission
    }{
        "Fetch page",
        submissions,
    }

	r.HTML(w, http.StatusOK, "admin/fetch", p)
}
func fetchDetailHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
    type Submission struct {
        Title string
        Permalink string
        RedditID string  `db:"reddit_id"`
        URL string
        Selftext string
    }
	submission := Submission{}
    err := conn.Get(&submission, "select title, permalink, reddit_id, url, selftext from submissions where reddit_id = ?", params.ByName("id"))
	if err != nil {
		logger.Println(err)
	}
    p := struct {
        Title string
        Sub Submission
    }{
    "detail page",
    submission,
    }
	r.HTML(w, http.StatusOK, "admin/fetch/detail", p)
}
func fetchJob(config Config) {
	ticker := time.Tick(1 * time.Minute)
	for range ticker {
		logger.Println("starting fetch job")

		var haves []string
		err := conn.Select(&haves, "select reddit_id from submissions")
		if err != nil {
			logger.Println(err)
		}
		fetch_sql, err := conn.Prepare("insert into submissions (url, reddit_id, permalink, title, submitted_by ,selftext, date_created) values ( ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			logger.Println(err)
		}

		ua := fmt.Sprintf("web:Meppit: /r/meppit -maps for reddit-:v %02f (by /u/anoland)", config.Version)

		ggg, err := geddit.NewLoginSession(config.Reddituser, config.Redditpass, ua)
        logger.Printf("%#v", ggg)
        if err != nil {
            logger.Printf("login session error: ", err)
        }
        subreddit := "meppit"
        if debug {
            subreddit = "meppitdev"
        }
		submissions, err := ggg.SubredditSubmissions(subreddit)
        if err != nil {
            logger.Println(err)
        }
		count := 0
		for _, sub := range submissions {
			_, ok := idExists(haves, sub.ID)
			if !ok {
				_, err := fetch_sql.Exec(sub.URL, sub.ID, sub.Permalink, sub.Title, sub.Author, sub.Selftext, sub.DateCreated)
				count = count + 1
				if err != nil {
					logger.Println(err)
				}
			}

		}
		// todo: keep track of last sucessful fetch time

		logger.Println("finished fetch job.")
		logger.Printf("processed %d submissions", count)
	}
}


func idExists(haystack []string, needle string) (int, bool) {
	sort.Strings(haystack)
	l := len(haystack)
	i := sort.SearchStrings(haystack, needle)
	if i < l && haystack[i] == needle {
		return i, true
	}
	return l, false

}
