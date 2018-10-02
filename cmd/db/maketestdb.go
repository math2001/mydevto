package main

import (
	"log"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test/testdb"
)

func main() {
	log.Println("Connecting to the database...")
	db := db.DB()
	log.Println("Creating the schema for the database...")
	_, err := db.Exec(`
	CREATE TABLE users (
		id        SERIAL,
		token     VARCHAR(255) NOT NULL DEFAULT '',
		service   VARCHAR(1024) NOT NULL,
		email     VARCHAR(255) NOT NULL,
		username  VARCHAR(255) NOT NULL,
		avatar    VARCHAR(255) NOT NULL DEFAULT '',
		name      VARCHAR(255) NOT NULL DEFAULT '',
		bio       VARCHAR(255) NOT NULL DEFAULT '',
		url       VARCHAR(255) NOT NULL DEFAULT '',
		location  VARCHAR(255) NOT NULL DEFAULT '',
		updated   TIMESTAMPTZ DEFAULT  now(),
		PRIMARY KEY (id),
		UNIQUE (email, service)
	);

	CREATE TABLE posts (
		id       SERIAL,
		userid   INTEGER,
		title    VARCHAR(255) NOT NULL,
		content  TEXT NOT NULL,
		written  TIMESTAMPTZ DEFAULT NOW(),
		updated  TIMESTAMPTZ DEFAULT NOW(),
		PRIMARY KEY (id),
		FOREIGN KEY (userid) REFERENCES users(id)
	);

	CREATE TABLE comments (
		id       SERIAL,
		userid   INTEGER NOT NULL,
		postid   INTEGER NOT NULL,
		content  TEXT NOT NULL,
		written  TIMESTAMPTZ DEFAULT NOW(),
		updated  TIMESTAMPTZ DEFAULT NOW(),
		PRIMARY KEY (id),
		FOREIGN KEY (postid) REFERENCES posts(id),
		FOREIGN KEY (userid) REFERENCES users(id)
	);`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Populating the database...")
	testdb.Populate()
	log.Println("Done.")
}
