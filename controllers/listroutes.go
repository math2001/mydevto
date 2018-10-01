package controllers

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/resp"
	"github.com/math2001/mydevto/services/uli"
)

// ListRoutes is a little utility that lists every routes that router (or
// therefore subrouter) supports. This is useful for devlopment
type ListRoutes struct {
	Router *mux.Router
}

func (lr ListRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var routes = make(map[string]string)
	err := lr.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
		uli.Printf(r, "Errored walking routes: %s", err)
		resp.InternalError(w, r)
		return
	}
	resp.Encode(w, r, routes)
}
