package main

import (
	"log"
	"net/http"
	"os"
	"time"

	// initiate the drivers for postgresql
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/math2001/mydevto/app"
	"github.com/math2001/mydevto/db"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set for this server to run")
	}

	dbconn, err := db.Open(db.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     logins["db"].username,
		Password: logins["db"].password,
		DBName:   "mydevto",
	})
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte(logins["session"].password))

	log.Printf("Connected to the database")

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Handle("/", http.FileServer(http.Dir("public")))
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
