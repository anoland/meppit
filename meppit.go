package main

import (
	//_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
    "os"
    "github.com/anoland/meppit/config"

	"github.com/anoland/geddit"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
)
var (
    logger *log.Logger
    rndr      *render.Render
)

func main() {

	logger := log.New(os.Stdout, "", log.Lshortfile)

    cfg, err := config.ParseConfig("config.json"); if err != nil {
        log.Fatal(err)
    }

	rndr = render.New(render.Options{
		Directory:     "templates",
		Layout:        "layout",
		IsDevelopment: cfg.Debug,
	})


	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/redirect", redirectHandler)
	router.GET("/admin/fetch/", fetchHandler)
	router.GET("/admin/fetch/detail/:id", fetchDetailHandler)

	listen := fmt.Sprintf("%s:%s", cfg.Host.ListenAddr, cfg.Host.HttpListenPort)
	logger.Println("Starting server on: ", listen)
	n := negroni.Classic()
	n.UseHandler(router)
    n.Run(listen)
}

func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p := struct {
		Title string
		Body  string
	}{
		"Index page",
		"This is the body",
	}
	rndr.HTML(w, http.StatusOK, "index", p)
}

func loginHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

        p := struct {
            Title string
            Body  string
            URL string
        }{
            "Login page",
            "token is valid",
            "" ,
        }
    if !oauthToken.Valid() {

        logger.Println("invalid token")
        //optionDuration := oauth2.SetAuthURLParam("duration", "permanent")
        loginurl := oauthCfg.AuthCodeURL("coding") //,  optionDuration)

        p = struct {
            Title string
            Body  string
            URL   string
        }{
            "Login page",
            "This where you login",
            loginurl,
        }

    }
    rndr.HTML(w, http.StatusOK, "login", p)
}

func redirectHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    
    code := req.FormValue("code")
    var err error
    oauthToken, err = oauthCfg.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Fatal(err)
    }

    client := oauthCfg.Client(oauth2.NoContext, oauthToken)

    me := geddit.Redditor{} 

    resp, err := client.Get("https://oauth.reddit.com/api/v1/me.json")
    body,_ := ioutil.ReadAll(resp.Body)
    raw_token := fmt.Sprintf("%#v",oauthToken)

    json.Unmarshal(body, &me)
    defer resp.Body.Close()

    if err != nil {
        log.Fatal(err)
    }
    insert_sql := "insert into tokens (username, token, expires, raw_token) values (?, ?, ?, ?)"
    conn.MustExec(insert_sql, me.Name, oauthToken.AccessToken, oauthToken.Expiry.Format(time.RFC3339), raw_token);

	p := struct {
		Title  string
		Body   string
        Me  geddit.Redditor
	}{
		"redirect page",
		"These are the params:",
		me,
	}
	rndr.HTML(w, http.StatusOK, "redirect", p)
}

func fetchHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	type Submission struct {
		Title     string
		Permalink string
		RedditID  string `db:"reddit_id"`
		URL       string
	}
	submissions := []Submission{}
	err := conn.Select(&submissions, "select title, permalink, reddit_id, url from submissions")
	if err != nil {
		logger.Println(err)
	}
	p := struct {
		Title       string
		Submissions []Submission
	}{
		"Fetch page",
		submissions,
	}

	rndr.HTML(w, http.StatusOK, "admin/fetch", p)
}
func fetchDetailHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	type Submission struct {
		Title     string
		Permalink string
		RedditID  string `db:"reddit_id"`
		URL       string
		Selftext  string
	}
	submission := Submission{}
	err := conn.Get(&submission, "select title, permalink, reddit_id, url, selftext from submissions where reddit_id = ?", params.ByName("id"))
	if err != nil {
		logger.Println(err)
	}
	p := struct {
		Title string
		Sub   Submission
	}{
		"detail page",
		submission,
	}
	rndr.HTML(w, http.StatusOK, "admin/fetch/detail", p)
}
