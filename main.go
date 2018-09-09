package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/math2001/mydevto/db"
)

func main() {
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
	fmt.Printf("Connected to the database! %v\n", dbconn)
}
