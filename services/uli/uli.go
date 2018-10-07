// Package uli provides an interface that automatically logs the request's
// identity as a prefix
package uli

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/math2001/mydevto/services/buildinfos"
)

type key int

const (
	logsdir    = "./logs"
	ctxRequest = key(iota)
)

var logger *log.Logger

func init() {
	// TODO: fix up the logs during testing. Please.
	if buildinfos.Testing {
		f, err := os.OpenFile(os.TempDir()+"/mydevto.logs",
			os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("uli: couldn't create temporary log file for testing: %s", err)
		}
		logger = log.New(f, "", log.LstdFlags)
		logger.Printf("\n---\n")
		return
	}
	f := createfile(buildinfos.V)
	// writes both to stdout and the file
	logger = log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags)
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
func Printf(ctx context.Context, format string, a ...interface{}) {
	r, ok := ctx.Value(ctxRequest).(*http.Request)
	if !ok {
		// there is no request in the context
		logger.Printf(format, a...)
		return
	}
	logger.Printf("%s %s %s", r.RemoteAddr, r.RequestURI, fmt.Sprintf(format, a...))
}

// Security display a warning header indicating that the next log is could be
// about a security flaw
func Security(ctx context.Context) {
	Printf(ctx, "POTENTIAL SECURITY FLAW")
}

// Middleware is the middleware that mux will use
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxRequest, r)
		start := time.Now()
		Printf(ctx, "Handling")
		next.ServeHTTP(w, r.WithContext(ctx))
		Printf(ctx, "Finished after %.2fs", time.Now().Sub(start).Seconds())
	})
}
