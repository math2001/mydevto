package router

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/api/posts"
	"github.com/math2001/mydevto/api/users"
	"github.com/math2001/mydevto/services/uli"
)

// This is the homepage, the only page the user is going to load directly
// quick and dirty. This is going to be managed by the views in the future.
// It breaks testing though. I enable it when I need it (to create JWT).
func index() http.HandlerFunc {
	githubid := os.Getenv("GITHUBID")
	if githubid == "" {
		log.Fatalf("$GITHUBID isn't defined.")
	}
	t, err := template.ParseFiles("web/index.tmpl")
	if err != nil {
		log.Fatalf("Errored parsing index.html: %s", err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, map[string]string{"GithubID": githubid}); err != nil {
		log.Fatalf("Errored executing template: %s", err)
	}
	var html = b.String()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	}
}

// Router creates the router for the application, delegating it to the
// controllers
func Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(uli.Middleware)
	router.StrictSlash(true)

	// router.HandleFunc("/", index())
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	manageAPI(router.PathPrefix("/api").Subrouter())
	return router
}

func manageAPI(r *mux.Router) {
	r.Handle("/", api.ListRoutes(r)).Methods("GET")
	posts.Manage(r.PathPrefix("/posts/").Subrouter())
	users.Manage(r.PathPrefix("/users/").Subrouter())
}
