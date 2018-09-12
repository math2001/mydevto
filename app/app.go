package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/math2001/mydevto/db"
)

const (
	sessionauth   = "authentication"
	servicegithub = iota
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
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	URL      string `json:"name"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	// Token    string `json:"-"`
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

// Init adds handlers to the router and initiates different stuff
func Init(router *mux.Router, d *db.Conn, s *sessions.FilesystemStore) {
	dbconn = d
	store = s
	r = router
	r.HandleFunc("/", home)
	r.HandleFunc("/posts", posts)
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
