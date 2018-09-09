package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/db"
)

var dbconn *db.Conn
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

func writeErr(w http.ResponseWriter, r *http.Request, msg string, code int) {
	fmt.Println("error", msg)
	enc(w, r, map[string]string{
		"type":    "error",
		"message": msg,
	})
	w.WriteHeader(code)
}

func internalErr(w http.ResponseWriter, r *http.Request) {
	writeErr(w, r, "Internal error", http.StatusInternalServerError)
}

// enc writes the object to the page, formatting according the User-Agent
func enc(w http.ResponseWriter, r *http.Request, e interface{}) error {
	enc := json.NewEncoder(w)
	if r.UserAgent() != "fetch.js" {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(e)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hum... This is the api. Just JSON.")
}

// Init adds handlers to the router
func Init(r *mux.Router, dbconnlocal *db.Conn) {
	dbconn = dbconnlocal
	r.HandleFunc("/", home)
	r.HandleFunc("/posts", posts)
}
