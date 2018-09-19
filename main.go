package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	// initiate the drivers for postgresql
	_ "github.com/lib/pq"
	"github.com/math2001/mydevto/app"
	"github.com/math2001/mydevto/db"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dblogin := os.Getenv("DBLOGIN")
	if dblogin == "" {
		log.Fatal("$DBLOGIN must be set")
	}
	dbconn, err := db.Open(db.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     dblogin,
		Password: os.Getenv("DBPASSWORD"),
		DBName:   "mydevto",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to the database")

	sessionkey := os.Getenv("SESSIONKEY")
	if sessionkey == "" {
		log.Fatal("$SESSIONKEY must be set")
	}

	store := sessions.NewFilesystemStore("", []byte(sessionkey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 30 days,
		HttpOnly: true,
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", app.Index)
	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	app.Init(r.PathPrefix("/api").Subrouter(), &dbconn, store)

	log.Printf("Running on :%s", port)

	server := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
