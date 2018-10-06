package db

import (
	"time"

	"github.com/fatih/structs"
	"github.com/gbrlsnchs/jwt"
)

// User represents user data
type User struct {
	ID       int       `json:"-"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Name     string    `json:"name,omitempty"`
	URL      string    `json:"url,omitempty"`
	Service  string    `json:"-"`
	Email    string    `json:"email,omitempty"`
	Location string    `json:"location,omitempty"`
	Bio      string    `json:"bio,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}

// JWTToken stores the user information in a web token
type JWTToken struct {
	*jwt.JWT
	*User
}

// Equals check if the non-zero fields in u are the same in o. This means that
// o might have extra fields, it's can still return true. In other words,
// the second struct is allowed to have more field than the first one.
func (u User) Equals(o User) bool {
	of := structs.New(o)
	for _, field := range structs.Fields(u) {
		// if the field isn't exported, or it hasn't been set, just ignore it
		if !field.IsExported() || field.IsZero() {
			continue
		}
		if field.Value() != of.Field(field.Name()).Value() {
			return false
		}
	}
	return true
}

// Post represents a post data
type Post struct {
	ID      int       `json:"-"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
	Written time.Time `json:"written"`
	Content string    `json:"content"`
	User    User      `json:"user"`
}

// Equals performs a similar check than user.Equals
func (u Post) Equals(o Post) bool {
	of := structs.New(o)
	for _, field := range structs.Fields(u) {
		// if the field isn't exported, or it hasn't been set, just ignore it
		if !field.IsExported() || field.IsZero() {
			continue
		}
		if field.Name() == "User" {
			// don't need type assertion, because the field User has to be User
			if !u.User.Equals(field.Value().(User)) {
				return false
			}
		} else if field.Value() != of.Field(field.Name()).Value() {
			return false
		}
	}
	return true
}
