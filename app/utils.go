package app

import (
	"encoding/json"
	"net/http"
)

func writeErr(w http.ResponseWriter, r *http.Request, msg string, code int) {
	enc(w, r, map[string]string{
		"type":    "error",
		"message": msg,
	})
	w.WriteHeader(code)
}

func internalErr(w http.ResponseWriter, r *http.Request) {
	writeErr(w, r, "Internal error", http.StatusInternalServerError)
}

// enc writes the object to the page, formatting according the User-Agent
func enc(w http.ResponseWriter, r *http.Request, e interface{}) error {
	enc := json.NewEncoder(w)
	if r.UserAgent() != "fetch.js" {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(e)
}
