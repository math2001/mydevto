package posts

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/math2001/mydevto/controllers/users"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/services/resp"
	"github.com/math2001/mydevto/services/uli"
)

const maxTitleLength = 255

func invaliddata(w http.ResponseWriter, r *http.Request) {
	resp.Error(w, r, http.StatusBadRequest, "Invalid request form data")
}

// write creates or updates a posts, depending on whether id is given
func write(w http.ResponseWriter, r *http.Request) {
	user := users.Current(r)
	if user == nil {
		resp.RequestLogin(w, r)
		return
	}
	userid := user.ID
	if err := r.ParseForm(); err != nil {
		uli.Printf(r, "Could not parse form data: %s", err)
		resp.Error(w, r, http.StatusBadRequest,
			"Could not parse form data: %s", err)
		return
	}
	// we don't give too much information about the error to the client, it's
	// probably because they're trying to change stuff manually.
	// To debug, just look at the logs
	var (
		id  int // remember, the default value is 0
		err error
	)
	// id (and therefore idstring) represent the *post* id. If it's given,
	// it means we update a post, otherwise, we just insert a new one
	idstring := r.PostForm.Get("id")
	if id, err = strconv.Atoi(idstring); idstring != "" && err != nil {
		uli.Printf(r, "Invalid id field: %q", idstring)
		invaliddata(w, r)
		return
	}
	title := r.PostForm.Get("title")
	if title == "" {
		uli.Printf(r, "no title given")
		invaliddata(w, r)
		return
	}
	if len(title) > maxTitleLength {
		uli.Printf(r, "Title length too long: %q (max is %d)", title,
			maxTitleLength)
		invaliddata(w, r)
	}
	content := r.PostForm.Get("content")
	if content == "" {
		uli.Printf(r, "no content given")
		invaliddata(w, r)
		return
	}

	var newid int
	if id == 0 {
		err = db.DB().QueryRow(`
			INSERT INTO posts (
				userid,
				title,
				content
			) VALUES ( $1, $2, $3 )
			RETURNING (id)
		`, userid, title, content).Scan(&newid)
		if err != nil {
			uli.Printf(r, "could not insert post new post: %s", err)
			// I sure hope no one will have to parse through that, but I'm sure
			// it'll be very helpful when this error occurs
			uli.Printf(r, "userid: %d title: %q content: %q", userid, title,
				content)
			resp.InternalError(w, r)
			return
		}
		resp.Encode(w, r, map[string]interface{}{
			"type":    "success",
			"message": "post successfully inserted",
			"id":      newid,
		}, http.StatusOK)
		return
	}
	err = db.DB().QueryRow(`
		UPDATE posts SET title=$1, content=$2, updated=NOW()
		WHERE id=$3 AND userid=$4
		RETURNING (id)
	`, title, content, id, userid).Scan(&newid)
	if err == sql.ErrNoRows {
		uli.Security(r)
		uli.Printf(r, "invalid combination postid (%d) and userid (%d)", id, userid)
		resp.InternalError(w, r)
		return
	} else if err != nil {
		uli.Printf(r, "could not update post: %s", err)
		resp.InternalError(w, r)
		return
	}
	// newid should be the same as id. I implement this behaviour because I'm
	// not sure about the behaviour of SERIAL in postgres (04.10.2018), which
	// is the type of the id field
	if newid != id {
		uli.Security(r)
		uli.Printf(r, "post id (%d) is different from the returned id (%d)",
			id, newid)
		resp.InternalError(w, r)
		return
	}
	resp.Success(w, r, "post updated successfully")
}
