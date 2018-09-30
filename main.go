package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
	"github.com/math2001/mydevto/controllers/posts"
	"github.com/math2001/mydevto/controllers/users"
	"github.com/math2001/mydevto/services/db"
)

var router *mux.Router

// This is the homepage, the only page the user is going to load directly
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

func initAPI(r *mux.Router) {
	r.Handle("/", controllers.ListRoutes{Router: r}).Methods("GET")
	posts.Manage(r.PathPrefix("/posts/").Subrouter())
	users.Manage(r.PathPrefix("/users/").Subrouter())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router = mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", index())
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	initAPI(router.PathPrefix("/api").Subrouter())

	db.Init()
	log.Printf("Running on :%s", port)

	server := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, router),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
