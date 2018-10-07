package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/services/uli"
)

const (
	// SessionAuth is the session name for sessions.FilesystemStore
	SessionAuth = "authentication"
	// ServiceGithub is the expected name that will be given as the oauth
	// callback
	ServiceGithub = "github"
	// JWT is the cookie name to store the token in
	JWT = "jwt"
)

// HTTPClient is the only http client that is allowed to be used by any part of
// this application. This is used to make sure that there is a timeout (there
// isn't any on the default one)
var HTTPClient = http.Client{
	Timeout: 20 * time.Second,
}

// ListRoutes is a little utility that lists every routes that router (or
// subrouter) handles. This is useful for devlopment
func ListRoutes(router *mux.Router) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var routes = make(map[string]string)
		err := router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			m, err := route.GetMethods()
			// this check if absolutely gross, but mux doesn't export their error,
			// so I don't really have a choice (Tuesday 25 September 2018)
			if err != nil && err.Error() != "mux: route doesn't have methods" {
				return err
			}
			if len(m) == 0 {
				routes[t] = "*"
			} else {
				routes[t] = strings.Join(m, ", ")
			}
			return nil
		})
		if err != nil {
			uli.Printf(ctx, "Errored walking routes: %s", err)
			InternalError(w, r)
			return
		}
		Encode(w, r, routes, http.StatusOK)
	}
}
