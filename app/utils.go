package app

import (
	"encoding/json"
	"net/http"
)

func writeMsg(w http.ResponseWriter, r *http.Request, kind string, msg string, code int) {
	w.WriteHeader(code)
	enc(w, r, map[string]string{
		"type":    kind,
		"message": msg,
	})
}

func writeErr(w http.ResponseWriter, r *http.Request, msg string, code int) {
	writeMsg(w, r, "error", msg, code)
}

func writeSuc(w http.ResponseWriter, r *http.Request, msg string) {
	writeMsg(w, r, "error", msg, http.StatusOK)
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
