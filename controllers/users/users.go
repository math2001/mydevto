package users

import (
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
)

// Manage is delegated the charges of mapping routes to functions by the main
// package
func Manage(r *mux.Router) {
	r.Handle("/", controllers.ListRoutes{Router: r}).Methods("GET")
	// in there documentation, github say they send a POST request, but they
	// actually send a GET... :(
	r.HandleFunc("/auth", auth).Methods("GET")
}
