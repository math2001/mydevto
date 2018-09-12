package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Masterminds/squirrel"
)

func posts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		postsGet(w, r)
	} else if r.Method == "POST" {
		postsPost(w, r)
	}
}

func postsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Errored parsing form data @ postsPost: %s", err)
		writeErr(w, r, fmt.Sprintf("Invalid form: %s", err), http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	if title == "" {
		log.Printf("No title provided @ postsPost")
		writeErr(w, r, "`title` POST field missing", http.StatusBadRequest)
		return
	}
	content := r.PostForm.Get("content")
	if content == "" {
		log.Printf("No content provided @ postsPost")
		writeErr(w, r, "`content` POST field missing", http.StatusBadRequest)
		return
	}
	session, err := store.Get(r, sessionauth)
	if err != nil {
		log.Printf("Errored getting session from store @ postsPost: %s", err)
		writeErr(w, r, "Couldn't load session. Try again.", http.StatusBadRequest)
		return
	}
	userid, ok := session.Values["userid"]
	if !ok {
		log.Printf("POST request unauthenticated @ postsPost.")
		writeErr(w, r, "Authenticate first", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Insert new post")
	fmt.Fprintln(w, "Title:", title)
	fmt.Fprintln(w, "Content:", content)
	fmt.Fprintln(w, "User id:", userid)
}

// lists posts, according to the parameters in the URL
func postsGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var (
		limit  int
		userid int
		err    error
	)
	if lim := q.Get("limit"); lim != "" {
		limit, err = strconv.Atoi(lim)
		if err != nil {
			log.Printf("Invalid limit @ postsGet: %s", err)
			writeErr(w, r, "Invalid `limit`. Should be a number", http.StatusBadRequest)
			return
		}
	}
	if ui := q.Get("userid"); ui != "" {
		userid, err = strconv.Atoi(ui)
		if err != nil {
			log.Printf("Invalid userid @ postsGet: %s", err)
			writeErr(w, r, "Invalid `userid`. Should be a number", http.StatusBadRequest)
			return
		}
	}
	req := psql.Select("p.title, p.written, u.username, u.avatar").
		From("posts p").
		LeftJoin("users u ON p.userid = u.id").
		OrderBy("written DESC")
	if limit > 0 {
		req = req.Limit(uint64(limit))
	}
	if userid != 0 {
		req = req.Where(squirrel.Eq{"userid": userid})
	}
	sql, args, err := req.ToSql()
	if err != nil {
		log.Printf("Errored building sql request @ postIndex: %s", err)
		internalErr(w, r)
		return
	}
	rows, err := dbconn.DB.Query(sql, args...)
	if err != nil {
		log.Printf("Errored querying @ postIndex: %s", err)
		internalErr(w, r)
		return
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title, &post.Updated, &post.User.Username, &post.User.Avatar); err != nil {
			log.Printf("Errored scanning @ postIndex: %s", err)
			internalErr(w, r)
			return
		}
		posts = append(posts, post)
	}
	if err := enc(w, r, posts); err != nil {
		log.Printf("Errored encoding @ postIndex: %s", err)
		internalErr(w, r)
		return
	}
}
