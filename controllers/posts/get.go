package posts

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
	"github.com/math2001/mydevto/resp"
	"github.com/math2001/mydevto/services/db"
)

// get gets a post by id
func get(w http.ResponseWriter, r *http.Request) {
	idstring := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstring)
	if err != nil {
		log.Printf("Couldn't convert id %q to integer: %s", idstring, err)
		resp.Error(w, r, fmt.Sprintf("Couldn't convert id %q to integer", idstring),
			http.StatusBadRequest)
		return
	}
	p := controllers.Post{}
	u := &p.User
	err = db.DB().QueryRow(`
	SELECT p.title, p.content, p.written, p.updated,
		   u.username, u.bio, u.url, u.avatar, u.name
	FROM posts AS p
	JOIN users AS u
	ON u.id = p.userid
	WHERE p.id=$1
	`, id).Scan(&p.Title, &p.Content, &p.Written, &p.Updated, &u.Username,
		&u.Bio, &u.URL, &u.Avatar, &u.Name)

	if err == sql.ErrNoRows {
		log.Printf("Couldn't find post with id %d", id)
		resp.Error(w, r, fmt.Sprintf("No post found with id %d", id),
			http.StatusBadRequest)
	} else if err != nil {
		log.Printf("Errored querying post from id %d: %s", id, err)
		resp.InternalError(w, r)
	} else {
		resp.Encode(w, r, p)
	}
}
