package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers/posts"
	"github.com/math2001/mydevto/controllers/users"
	"github.com/math2001/mydevto/resp"
	"github.com/math2001/mydevto/services/db"
)

var router *mux.Router
var apirouter *mux.Router

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
	var routes = make(map[string]string)
	err := apirouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		// this check if absolutely gross, but mux doesn't export their error,
		// so I don't really have a choice (Tuesday 25 September 2018)
		if err != nil && err.Error() != "mux: route doesn't have methods" {
			return err
		}
		if len(m) == 0 {
			routes[t] = "*"
		} else {
			routes[t] = strings.Join(m, ", ")
		}
		return nil
	})
	if err != nil {
		log.Printf("Errored walking routes: %s", err)
		resp.InternalError(w, r)
		return
	}
	resp.Encode(w, r, routes)
}

func initAPI() {
	apirouter.HandleFunc("/", apiIndex).Methods("GET")
	posts.Manage(apirouter.PathPrefix("/posts").Subrouter())
	users.Manage(apirouter.PathPrefix("/users").Subrouter())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router = mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", index)
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	apirouter = router.PathPrefix("/api").Subrouter()
	initAPI()

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
