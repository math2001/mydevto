// Package resp provides basic functions to write message as an http response
package resp

import (
	"encoding/json"
	"net/http"
)

// Message writes a response in the form { "type": <str>, "message": <str> }
// It is mainly a util function for the Error and Success functions
func Message(w http.ResponseWriter, r *http.Request, kind string, msg string, code int) error {
	w.WriteHeader(code)
	return Encode(w, r, map[string]string{
		"type":    kind,
		"message": msg,
	})
}

// Error writes a message of type "error" with "message" msg. It is typically
// used to give more information about an error to the user. This is when the
// information isn't sensitive, otherwise InternalError is used.
func Error(w http.ResponseWriter, r *http.Request, msg string, code int) error {
	return Message(w, r, "error", msg, code)
}

// Success writes a message of type "success" with "message" msg. It is
// typically used to confirm that a post request has been successful
func Success(w http.ResponseWriter, r *http.Request, msg string) error {
	return Message(w, r, "success", msg, http.StatusOK)
}

// InternalError writes a static object:
// {"kind": "error", "message": "Internal error"}
// This is used a lot, when the error that caused the failure could be
// sensitive, and we don't want to give too much information
func InternalError(w http.ResponseWriter, r *http.Request) error {
	return Error(w, r, "Internal error", http.StatusInternalServerError)
}

// Encode writes the object to the page, formatting according the User-Agent
func Encode(w http.ResponseWriter, r *http.Request, e interface{}) error {
	enc := json.NewEncoder(w)
	if r.UserAgent() != "js" {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(e)
}
