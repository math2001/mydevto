package users

import (
	"github.com/gorilla/mux"
)

// Manage is delegated the charges of mapping routes to functions by the main
// package
func Manage(r *mux.Router) {
	r.HandleFunc("/auth", auth)
}
