package users

import (
	"net/http"

	"github.com/math2001/mydevto/services/resp"
)

// current serves the information about the currently logged in user
func current(w http.ResponseWriter, r *http.Request) {
	u := Current(r)
	if u == nil {
		resp.Error(w, r, http.StatusForbidden, "please log in")
		return
	}
	resp.Encode(w, r, u)
}
