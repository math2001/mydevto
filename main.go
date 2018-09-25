package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers/posts"
	// initiate the drivers for postgresql
	_ "github.com/lib/pq"
)

// This is the homepage, the only page the user is going to load directly
func index(w http.ResponseWriter, r *http.Request) {
	// TODO: this is dumb. When we return error, the user doesn't want json
	// (which is what we write when we do internalErr)
	githubid := os.Getenv("GITHUBID")
	if githubid == "" {
		log.Printf("$GITHUBID isn't defined.")
		http.Error(w, "Invalid configuration: GITHUBID isn't defined",
			http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("web/index.tmpl")
	if err != nil {
		log.Printf("Errored parsing index.html: %s", err)
		http.Error(w, "Couldn't parse index template.", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, map[string]string{"GithubID": githubid}); err != nil {
		log.Printf("Errored executing template: %s", err)
		http.Error(w, "Couldn't execute template.", http.StatusInternalServerError)
		return
	}
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func initAPI(r *mux.Router) {
	r.HandleFunc("/", apiIndex)
	posts.Manage(r.PathPrefix("/posts").Subrouter())
	users.Manage(r.PathPrefix("/users").Subrouter())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", index)
	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	initAPI(r.PathPrefix("/api").Subrouter())

	log.Printf("Running on :%s", port)

	server := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
