package controllers

import (
	"net/http"
	"time"
)

const (
	// SessionAuth is the session name for sessions.FilesystemStore
	SessionAuth = "authentication"
	// ServiceGithub is the expected name that will be given as the oauth
	// callback
	ServiceGithub = "github"
)

// HTTPClient is the only http client that is allowed to be used by any part of
// this application. This is used to make sure that there is a timeout (there
// isn't any on the default one)
var HTTPClient = http.Client{
	Timeout: 20 * time.Second,
}

// Post represents a post data
type Post struct {
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Written string    `json:"written,omitempty"`
	Content string    `json:"content,omitempty"`
	User    User      `json:"user"`
}

// User represents user data
type User struct {
	ID       string `json:"-"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
}
