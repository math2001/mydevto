package sessions

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
)

// Creates a one instance of a file system store, and gives it to anyone who
// wants it

var store *sessions.FilesystemStore

func init() {
	sessionkey := os.Getenv("SESSIONKEY")
	if sessionkey == "" {
		log.Fatal("$SESSIONKEY must be set")
	}
	store = sessions.NewFilesystemStore("./.local/sessions", []byte(sessionkey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
	}
}

// Store returns the store to use
func Store() *sessions.FilesystemStore {
	return store
}
