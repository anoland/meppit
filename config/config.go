package config

import (

    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "golang.org/x/oauth2"
	"github.com/jmoiron/sqlx"
)

var (
    logger = *log.New(os.Stdout, "", log.Lshortfile)
    conn  *sqlx.DB
)
type Config struct {
    Version  float64  `json:"version"`
    Debug  bool
    Host     host     `json:"host"`
    Database database `json:"database"`
    Reddit   reddit   `json:"reddit"`
    OauthCfg oauth2.Config  
}

type host struct {
    ListenAddr string `json:"listenaddr"`
    HttpListenPort string `json:"httplistenport"`
    HttpsListenPort string `json:"httpslistenport"`
}
type database struct {
    Dbhost string `json:"dbhost"`
    Dbuser string `json:"dbuser"`
    Dbpass string `json:"dbpass"`
    Dbname string `json:"dbname"`
}
type reddit struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    ClientKey string `json:"clientkey"`
    SecretKey string `json:"secretkey"`
    RedirectURL string `json:"redirecturl"`
}

var RedditEndpoint = oauth2.Endpoint{
    AuthURL:  "https://www.reddit.com/api/v1/authorize",
    TokenURL: "https://www.reddit.com/api/v1/access_token",
}



func ParseConfig(filename string) (*Config, error) {

    logger := log.New(os.Stdout, "", log.Lshortfile)
    c := &Config{}
    cf, err := ioutil.ReadFile(filename)
    if err != nil {
        return c,err
    }
    if err := json.Unmarshal(cf, &c); err != nil {
        return c,err 
    }

    logger.Println("Running version:", c.Version)

    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", c.Database.Dbuser, c.Database.Dbpass, c.Database.Dbhost, c.Database.Dbname)
    conn, _ = sqlx.Open("mysql", dsn)
    if err := conn.Ping(); err != nil {
        logger.Fatal(err)
    }

    c.OauthCfg = oauth2.Config{
        ClientID:     c.Reddit.ClientKey,
        ClientSecret: c.Reddit.SecretKey,
        Endpoint:     RedditEndpoint,
        RedirectURL:  c.Reddit.RedirectURL,
        Scopes:       []string{"read", "identity"},
    }
    app_env := os.Getenv("APP_ENVIRONMENT")
    logger.Println("Running in environment: ", app_env)
    if app_env == "development" {
        c.Debug = true
        logger.Println("debugging is on")
    } else {
        c.Debug = false
        logger.Println("debugging is off")
    }
    return c, nil

}
