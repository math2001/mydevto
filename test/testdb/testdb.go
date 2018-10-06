package testdb

import (
	"log"
	"time"

	"github.com/math2001/mydevto/services/db"
)

// Users are the users populating the test database
var Users = []db.User{
	{
		ID:       1,
		Username: "test1",
		Name:     "test1",
		Bio:      "I'm test1",
		Email:    "test1@tests.com",
		Service:  "fake",
		Updated:  time.Unix(0, 0),
	},
	{
		ID:       2,
		Username: "test2",
		Name:     "test2",
		Bio:      "I'm test2",
		Email:    "test2@tests.com",
		Service:  "fake",
		Updated:  time.Unix(0, 0),
	},
}

// Posts are hte post populating the test database
var Posts = []db.Post{
	{
		ID:      1,
		Title:   "First",
		Content: "The first post",
		Updated: time.Unix(0, 0),
		Written: time.Unix(0, 0),
		User:    Users[0],
	},
	{
		ID:      2,
		Title:   "Second",
		Content: "The second post",
		Updated: time.Unix(10, 0),
		Written: time.Unix(10, 0),
		User:    Users[0],
	},
	{
		ID:      3,
		Title:   "Third",
		Content: "The third post",
		Updated: time.Unix(20, 0),
		Written: time.Unix(20, 0),
		User:    Users[1],
	},
}

// Populate the test database
func Populate() {
	for _, u := range Users {
		_, err := db.DB().Exec(`INSERT INTO users (id, username, name, bio, email,
		service, updated) VALUES ($1, $2, $3, $4, $5, $6, $7)`, u.ID, u.Username, u.Name, u.Bio,
			u.Email, u.Service, u.Updated)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, p := range Posts {
		_, err := db.DB().Exec(`INSERT INTO posts (title, content, updated, written,
		userid) VALUES ($1, $2, $3, $4, $5)`, p.Title, p.Content, p.Written,
			p.Updated, p.User.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
