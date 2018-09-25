package post

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/math2001/goctrl/db"
)

// Posts is a controller
type Posts struct {
	DB    *db.DB
	Store sessions.FilesystemStore
}

// Index manages the /posts URL
func (p Posts) Index(w http.ResponseWriter, r *http.Request) {

}
