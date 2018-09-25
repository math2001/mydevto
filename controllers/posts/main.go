package post

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/math2001/goctrl/db"
)

// Post represents a post data
type Post struct {
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Written string    `json:"written,omitempty"`
	Content string    `json:"content,omitempty"`
	User    User      `json:"user"`
}

// Posts is a controller
type Posts struct {
	DB    *db.DB
	Store sessions.FilesystemStore
}

// Index manages the /posts URL
func (p Posts) Index(w http.ResponseWriter, r *http.Request) {

}
