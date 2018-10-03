package posts

import (
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/controllers"
)

// Manage defines which routes are triggering specific functions
func Manage(r *mux.Router) {
	r.Handle("/", controllers.ListRoutes{Router: r}).Methods("GET")
	r.HandleFunc("/get", get).Methods("GET")
	r.HandleFunc("/list", list).Methods("GET")
}
