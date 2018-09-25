package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// initiate the drivers for postgresql
	_ "github.com/lib/pq"
	"github.com/math2001/mydevto/app"
	"github.com/math2001/mydevto/db"
)

type Services struct{}

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

	services = map[string]interface{}{
		"db":    &dbconn,
		"store": store,
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", app.Index)
	r.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	app.Init(r.PathPrefix("/api").Subrouter(), services)

	log.Printf("Running on :%s", port)

	server := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
