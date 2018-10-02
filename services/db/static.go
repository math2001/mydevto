package db

import "time"

// Post represents a post data
type Post struct {
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Written time.Time `json:"written,omitempty"`
	Content string    `json:"content,omitempty"`
	User    User      `json:"user"`
}

// User represents user data
type User struct {
	ID       int       `json:"-"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Service  string    `json:"-"`
	Email    string    `json:"email"`
	Location string    `json:"location"`
	Bio      string    `json:"bio"`
	Updated  time.Time `json:"updated"`
}
