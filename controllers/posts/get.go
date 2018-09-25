package posts

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/resp"
)

// get gets a post by id
func get(w http.ResponseWriter, r *http.Request) {
	resp.Encode(w, r, mux.Vars(r))
}
