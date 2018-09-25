package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/math2001/mydevto/db"
)

const (
	sessionauth   = "authentication"
	servicegithub = "github"
)

var (
	dbconn     *db.Conn
	store      *sessions.FilesystemStore
	r          *mux.Router
	psql       = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	httpclient = http.Client{
		Timeout: 20 * time.Second,
	}
)

// User represents user data
type User struct {
	ID       string `json:"-"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
}

// Post represents a post data
type Post struct {
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Written string    `json:"written,omitempty"`
	Content string    `json:"content,omitempty"`
	User    User      `json:"user"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hum... This is the api. Just JSON.")
}

// Index manages the index page (the only page that is directly loaded by the user)
func Index(w http.ResponseWriter, r *http.Request) {
	// TODO: this is dumb. When we return error, the user doesn't want json
	// (which is what we write when we do internalErr)
	githubid := os.Getenv("GITHUBID")
	if githubid == "" {
		log.Printf("$GITHUBID isn't defined. Returning")
		internalErr(w, r)
		return
	}
	t, err := template.ParseFiles("web/index.tmpl")
	if err != nil {
		log.Printf("Errored parsing index.html: %s", err)
		internalErr(w, r)
		return
	}
	if err := t.Execute(w, map[string]string{"GithubID": githubid}); err != nil {
		log.Printf("Errored executing template: %s", err)
		internalErr(w, r)
		return
	}
}

// Init adds handlers to the router and initiates different stuff
func Init(router *mux.Router, services map[string]interface{}) {
	dbconn = services["db"].(*db.Conn)
	store = services["store"].(sessions.FileSystemStore)
	r = router
	r.HandleFunc("/", home)
	handlePosts(r.PathPrefix("/posts").Subrouter())
	r.HandleFunc("/users/{action}", users)
	r.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, sessionauth)
		if err != nil {
			log.Fatal(err)
		}
		session.Values["userid"] = 1
		session.Save(r, w)
		fmt.Fprintf(w, "Wrote userid: %d", session.Values["userid"])
	})
}
