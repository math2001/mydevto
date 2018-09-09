package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler handles every request made to the api
func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(vars)
	fmt.Fprintf(w, "Good!")
}

// Init adds handlers to the router
func Init(r *mux.Router) {
	r.HandleFunc("/", Handler)
}
