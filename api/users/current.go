package users

import (
	"net/http"

	"github.com/math2001/mydevto/api"
)

// current serves the information about the currently logged in user
func current(w http.ResponseWriter, r *http.Request) {
	u := Current(r.Context(), r)
	if u == nil {
		api.RequestLogin(w, r)
		return
	}
	api.Encode(w, r, u, http.StatusOK)
}
