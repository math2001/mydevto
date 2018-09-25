package users

import (
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
)

// Manage is delegated the charges of mapping routes to functions by the main
// package
func Manage(r *mux.Router) {
	r.Handle("/", controllers.ListRoutes{Router: r}).Methods("GET")
	r.HandleFunc("/auth", auth)
}
