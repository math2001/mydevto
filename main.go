package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/math2001/mydevto/router"

	// init services
	"github.com/math2001/mydevto/services/buildinfos"
	_ "github.com/math2001/mydevto/services/db"
)

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

func main() {
	log.Println("MyDevTo", buildinfos.V)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Printf("Running on :%s", port)

	server := &http.Server{
		Handler:      router.Router(),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
