package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func postsIndex(w http.ResponseWriter, r *http.Request) {
	rows, err := dbconn.DB.Query(`SELECT p.title, p.written, u.username, u.avatar FROM posts p LEFT JOIN users u ON p.userid = u.id ORDER BY written DESC LIMIT 10`)
	if err != nil {
		log.Printf("Errored querying @ postIndex: %s", err)
		internalError(w, r)
		return
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title, &post.Updated, &post.User.Username, &post.User.Avatar); err != nil {
			log.Printf("Errored scanning @ postIndex: %s", err)
			internalError(w, r)
			return
		}
		posts = append(posts, post)
	}
	if err := enc(w, r, posts); err != nil {
		log.Printf("Errored encoding @ postIndex: %s", err)
		internalError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func postsAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars)
}
