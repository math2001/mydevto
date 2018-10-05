package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/api/posts"
	"github.com/math2001/mydevto/api/users"
	"github.com/math2001/mydevto/services/uli"
)

// Router creates the router for the application, delegating it to the
// controllers
func Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(uli.Middleware)
	router.StrictSlash(true)

	// router.HandleFunc("/", index())
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	manageAPI(router.PathPrefix("/api").Subrouter())
	return router
}

func manageAPI(r *mux.Router) {
	r.Handle("/", api.ListRoutes{Router: r}).Methods("GET")
	posts.Manage(r.PathPrefix("/posts/").Subrouter())
	users.Manage(r.PathPrefix("/users/").Subrouter())
}
