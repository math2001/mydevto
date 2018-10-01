package users

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
	"github.com/math2001/mydevto/services/sess"
	"github.com/math2001/mydevto/services/uli"
	"github.com/mitchellh/mapstructure"
)

// Manage is delegated the charges of mapping routes to functions by the main
// package
func Manage(r *mux.Router) {
	r.Handle("/", controllers.ListRoutes{Router: r}).Methods("GET")
	// in there documentation, github say they send a POST request, but they
	// actually send a GET... :(
	r.HandleFunc("/auth", auth).Methods("GET")
	r.HandleFunc("/current", current).Methods("GET")
}

// Current returns the current user's information from the sessions. It returns
// nil if he isn't connected
func Current(r *http.Request) *controllers.User {
	session, err := sess.Store().Get(r, controllers.SessionAuth)
	if _, ok := err.(*os.PathError); ok {
		// we assume that it's a "file not found"
		// To do it properly, we'd need to use syscall, which, by definition,
		// is platform dependent.
		// Anyway, even if it isn't, the results the same: the user isn't
		// logged in
		return nil
	} else if err != nil {
		uli.Printf(r, "Errored getting auth session: %s", err)
		return nil
	}
	if len(session.Values) == 0 {
		// nothing has been found
		return nil
	}
	u := &controllers.User{}
	mapstructure.Decode(session.Values, u)
	return u
}
