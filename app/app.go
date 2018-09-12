package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/math2001/mydevto/db"
)

var dbconn *db.Conn
var store *sessions.CookieStore
var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// User represents a user data
type User struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Token    string `json:"-"`
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
func Init(r *mux.Router, d *db.Conn, s *sessions.CookieStore) {
	dbconn = d
	store = s
	r.HandleFunc("/", home)
	r.HandleFunc("/posts", posts)
}
