package controllers

import "time"

const (
	// SessionAuth is the session name for sessions.FilesystemStore
	SessionAuth = "authentication"
	// ServiceGithub is the expected name that will be given as the oauth
	// callback
	ServiceGithub = "github"
)

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
