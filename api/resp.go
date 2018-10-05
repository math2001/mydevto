package api

// provides basic functions to write JSON as an http response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/math2001/mydevto/services/uli"
)

// Message writes a response in the form { "type": <str>, "message": <str> }
// It is mainly a util function for the Error and Success functions
func Message(w http.ResponseWriter, r *http.Request, code int, kind string, format string, a ...interface{}) {
	Encode(w, r, map[string]string{
		"type":    kind,
		"message": fmt.Sprintf(format, a...),
	}, code)
}

// Error writes a message of type "error" with "message" msg. It is typically
// used to give more information about an error to the user. This is when the
// information isn't sensitive, otherwise InternalError is used.
func Error(w http.ResponseWriter, r *http.Request, code int, format string, a ...interface{}) {
	Message(w, r, code, "error", format, a...)
}

// Success writes a message of type "success" with "message" msg. It is
// typically used to confirm that a post request has been successful
func Success(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) {
	Message(w, r, http.StatusOK, "success", format, a...)
}

// RequestLogin is similar to InternalErro. It just writes that the user should
// login before trying to access this part of the website
func RequestLogin(w http.ResponseWriter, r *http.Request) {
	uli.Printf(r, "%s tried to claim access", r.RemoteAddr)
	Error(w, r, http.StatusForbidden, "please log in")
}

// InternalError writes a static object:
// {"kind": "error", "message": "Internal error"}
// This is used a lot, when the error that caused the failure could be
// sensitive, and we don't want to give too much information
func InternalError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusInternalServerError, "Internal error")
}

// Encode writes the object to the page, formatting according the User-Agent
func Encode(w http.ResponseWriter, r *http.Request, e interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	if r.UserAgent() != "js" {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(e); err != nil {
		log.Printf("Errored encoding object to ResponseWriter (writing manually): %s", err)
		// manually print out some JSON
		_, err = fmt.Fprintf(w, `{ "type": "error", "message": "Internal error" }`)
		if err != nil {
			// there's nothing we can do now...
			log.Printf("Errored writing to ResponseWriter: %s", err)
		}
	}
}
