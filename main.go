package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/math2001/mydevto/router"

	// init services
	"github.com/math2001/mydevto/services/buildinfos"
	_ "github.com/math2001/mydevto/services/db"
)

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
