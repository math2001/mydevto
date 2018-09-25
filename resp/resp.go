// Package resp provides basic functions to write message as an http response
package resp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Message writes a response in the form { "type": <str>, "message": <str> }
// It is mainly a util function for the Error and Success functions
func Message(w http.ResponseWriter, r *http.Request, code int, kind string, format string, a ...interface{}) error {
	w.WriteHeader(code)
	return Encode(w, r, map[string]string{
		"type":    kind,
		"message": fmt.Sprintf(format, a...),
	})
}

// Error writes a message of type "error" with "message" msg. It is typically
// used to give more information about an error to the user. This is when the
// information isn't sensitive, otherwise InternalError is used.
func Error(w http.ResponseWriter, r *http.Request, code int, format string, a ...interface{}) error {
	return Message(w, r, code, "error", format, a...)
}

// Success writes a message of type "success" with "message" msg. It is
// typically used to confirm that a post request has been successful
func Success(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) error {
	return Message(w, r, http.StatusOK, "success", format, a...)
}

// InternalError writes a static object:
// {"kind": "error", "message": "Internal error"}
// This is used a lot, when the error that caused the failure could be
// sensitive, and we don't want to give too much information
func InternalError(w http.ResponseWriter, r *http.Request) error {
	return Error(w, r, http.StatusInternalServerError, "Internal error")
}

// HandleError logs the error and writes internal error. The idea is to make
// how sensitive errors are handle easy to change (log isn't the best thing)
func HandleError(w http.ResponseWriter, r *http.Request, err error, format string, a ...interface{}) {
	log.Print(errors.Wrapf(err, fmt.Sprintf(format, a...)))
	InternalError(w, r)
}

// Encode writes the object to the page, formatting according the User-Agent
func Encode(w http.ResponseWriter, r *http.Request, e interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if r.UserAgent() != "js" {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(e)
}
