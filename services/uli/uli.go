// Package uli provides an interface that automatically logs the request's
// identity as a prefix
package uli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const logsdir = "./logs"

var logger *log.Logger

// Init inits the logger
func Init(version string) {
	f := createfile(version)
	// writes both to stdout and the file
	w := io.MultiWriter(os.Stdout, f)
	logger = log.New(w, "", log.LstdFlags)
}

// createdirs makes sure all the directory where the logs are going to be
// written exist
func createfile(version string) io.Writer {
	// creates logsdir if it doesn't exist
	if _, err := os.Stat(logsdir); os.IsNotExist(err) {
		if err := os.Mkdir(logsdir, 0755); err != nil {
			log.Fatalf("uli: couldn't create log directory %q: %s", logsdir, err)
		}
	}
	versiondir := path.Join(logsdir, version)
	if _, err := os.Stat(versiondir); os.IsNotExist(err) {
		if err := os.Mkdir(versiondir, 0755); err != nil {
			log.Fatalf("uli: couldn't create version directory %q: %s", versiondir, err)
		}
	}
	t := time.Now()
	name := fmt.Sprintf("%d-%d-%d.%d-%d-%d.log", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	p := path.Join(versiondir, name)
	f, err := os.Create(p)
	if err != nil {
		log.Fatalf("uli: couldn't create log file %q: %s", p, err)
	}
	return f
}

// Printf logs the message with the request identification at the beginning.
func Printf(r *http.Request, format string, a ...interface{}) {
	logger.Printf("%s %s", r.RemoteAddr, fmt.Sprintf(format, a...))
}

// Middleware is the middleware that mux will use
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Printf(r, "Start handling request...")
		next.ServeHTTP(w, r)
		Printf(r, "Finished handling request.")
	})
}
