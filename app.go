package app

import (
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

func initAPI(r *mux.Router) {
	r.HandleFunc("/", apiIndex)
	posts.Manage(r.PathPrefix("/posts").Subrouter())
	users.Manage(r.PathPrefix("/users").Subrouter())
}
