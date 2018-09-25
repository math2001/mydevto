package posts

import (
	"github.com/gorilla/mux"
)

// Manage defines which routes are triggering specific functions
func Manage(r *mux.Router) {
	r.HandleFunc("/get/{id}", get)
}
