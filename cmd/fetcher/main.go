package main

import (
    "fmt"
    "log"
    "time"
    "os"
    "sort"
    "github.com/anoland/meppit/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/anoland/geddit"
	"github.com/jmoiron/sqlx"
)
var (
    conn *sqlx.DB
    logger *log.Logger
)

func main() {

	logger = log.New(os.Stdout, "", log.Lshortfile)

    cfg, err := config.ParseConfig("config.json"); if err != nil {
        log.Fatal(err)
    }

	go submissionsJob(cfg)

}
func submissionsJob(c *config.Config) {
	ticker := time.Tick(15 * time.Minute)
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

		ua := fmt.Sprintf("web:Meppit: /r/meppit -maps for reddit-:v %02f (by /u/anoland)", c.Version)
        ggg := geddit.NewSession(ua)
		logger.Printf("%#v", ggg)
		if err != nil {
			logger.Printf("login session error: ", err)
		}
		subreddit := "meppit"
		if c.Debug {
			subreddit = "meppitdev"
		}
		options := geddit.ListingOptions{}
		submissions, err := ggg.SubredditSubmissions(subreddit, geddit.DefaultPopularity, options)
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

