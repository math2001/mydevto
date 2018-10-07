package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/services/uli"
)

// get gets a post by id
func get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idstring := r.URL.Query().Get("id")
	if idstring == "" {
		uli.Printf(ctx, "Given an empty idstring")
		api.Error(w, r, http.StatusBadRequest, "Invalid id. Empty strings not allowed")
		return
	}
	id, err := strconv.Atoi(idstring)
	if err != nil {
		uli.Printf(ctx, "Couldn't convert id %q to integer: %s", idstring, err)
		api.Error(w, r, http.StatusBadRequest,
			"Couldn't convert id %q to integer", idstring)
		return
	}
	p := db.Post{}
	u := &p.User
	err = db.QueryRowContext(ctx, `
	SELECT p.title, p.content, p.written, p.updated,
		   u.username, u.bio, u.url, u.avatar, u.name
	FROM posts AS p
	JOIN users AS u
	ON u.id = p.userid
	WHERE p.id=$1
	`, id).Scan(&p.Title, &p.Content, &p.Written, &p.Updated, &u.Username,
		&u.Bio, &u.URL, &u.Avatar, &u.Name)

	if err == sql.ErrNoRows {
		uli.Printf(ctx, "No post found with id %d", id)
		api.Error(w, r, http.StatusBadRequest, "No post found with id %d", id)
	} else if err != nil {
		uli.Printf(ctx, "Errored querying post from id %d: %s", id, err)
		api.InternalError(w, r)
	} else {
		api.Encode(w, r, p, http.StatusOK)
	}
}
