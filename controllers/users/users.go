package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/services/db"
	"github.com/pkg/errors"
)

func users(w http.ResponseWriter, r *http.Request) {
	action, ok := mux.Vars(r)["action"]
	if !ok {
		log.Printf("No action specified @ users")
		writeErr(w, r, "No action specified", http.StatusBadRequest)
		return
	}
	if action == "logout" {
		usersLogout(w, r)
	} else if action == "auth" {
		usersAuth(w, r)
	} else if action == "get" {
		usersGet(w, r)
	} else {
		log.Printf("Unknown action %q", action)
		writeErr(w, r, "Unkown action", http.StatusBadRequest)
		return
	}
}

// usersAuth is called by the service with 'service' as the name in the
// parameters, and finishes the auth flow. It's the 'callback'
func usersAuth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	service := query.Get("service")
	if service != "github" {
		log.Printf("Invalid service %q trying to authenticate @ usersAuth", service)
		writeErr(w, r, "Invalid service", http.StatusBadRequest)
		return
	}
	sessioncode := query.Get("code")
	if sessioncode == "" {
		log.Printf("No session code in URL parameter @ usersAuth#%q", service)
		internalErr(w, r)
		return
	}
	token, err := getToken(sessioncode)
	if err != nil {
		log.Printf("Errored getting token: %s", err)
		internalErr(w, r)
		return
	}
	user, err := retrieveUserInformation(token, servicegithub)
	if err != nil {
		log.Printf("Errored retrieving user information from token: %s", err)
		internalErr(w, r)
		return
	}
	id, err := saveUserInformation(token, servicegithub, user)
	if err != nil {
		log.Printf("Errored saving user information to database: %s", err)
		internalErr(w, r)
		return
	}
	session, err := store.Get(r, sessionauth)
	if err != nil {
		log.Printf("Errored getting authentication session @ usersAuth: %s", err)
		internalErr(w, r)
		return
	}
	session.Values["id"] = id

	session.Values["username"] = user.Username
	session.Values["avatar"] = user.Avatar
	session.Values["name"] = user.Name
	session.Values["bio"] = user.Bio
	session.Values["url"] = user.URL
	session.Values["email"] = user.Email
	session.Values["location"] = user.Location

	session.Save(r, w)
	fmt.Fprintf(w, "<script>window.close()</script>")
}

// removes the session cookie. Due to GitHub, we can't invalidate the token
// though
func usersLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionauth)
	if err != nil {
		log.Printf("Errored getting authentication session @ usersLogout: %s", err)
		internalErr(w, r)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	writeSuc(w, r, "Logged out")
}

// usersGet returns the information about the logged in user
func usersGet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionauth)
	if err != nil {
		log.Printf("Errored getting authentication session @ usersGet: %s", err)
		internalErr(w, r)
		return
	}
	if _, ok := session.Values["id"]; !ok {
		writeErr(w, r, "Please login", http.StatusBadRequest)
		return
	}
	enc(w, r, map[string]interface{}{
		"id":       session.Values["id"],
		"username": session.Values["username"],
		"avatar":   session.Values["avatar"],
		"name":     session.Values["name"],
		"bio":      session.Values["bio"],
		"url":      session.Values["url"],
		"email":    session.Values["email"],
		"location": session.Values["location"],
	})
}

func getToken(sessioncode string) (string, error) {
	params := url.Values{}
	githubid := os.Getenv("GITHUBID")
	if githubid == "" {
		return "", fmt.Errorf("$GITHUBID isn't defined. Aborting authentification")
	}
	githubsecret := os.Getenv("GITHUBSECRET")
	if githubsecret == "" {
		return "", fmt.Errorf("$GITHUBSECRET isn't defined. Aborting authentification")
	}
	params.Set("client_id", githubid)
	params.Set("client_secret", githubsecret)
	params.Set("code", sessioncode)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		strings.NewReader(params.Encode()))
	if err != nil {
		return "", errors.Wrapf(err, "errored building request getting token")
	}
	req.Header.Add("Accept", "application/json")
	res, err := httpclient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "errored doing request getting token")
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var result map[string]interface{}
	err = dec.Decode(&result)
	if err != nil {
		return "", errors.Wrapf(err, "errored decoding json getting token")
	}
	val, ok := result["access_token"]
	if !ok {
		return "", errors.New("no 'access_token' key in response getting token")
	}
	token, ok := val.(string)
	if !ok {
		return "", errors.New("errored asserting token as a string getting token")
	}
	return token, nil
}

func retrieveUserInformation(token string, service string) (User, error) {
	var user User
	if service == servicegithub {
		req, err := http.NewRequest("GET", "https://api.github.com/user?access_token="+token, nil)
		if err != nil {
			return user, errors.Wrapf(err, "errored creating request")
		}
		req.Header.Add("Accept", "application/json")
		res, err := httpclient.Do(req)
		if err != nil {
			return user, errors.Wrapf(err, "errored doing request")
		}
		if res.StatusCode != http.StatusOK {
			return user, errors.Errorf("invalid status code: %d", res.StatusCode)
		}
		defer res.Body.Close()
		dec := json.NewDecoder(res.Body)
		var result map[string]interface{}
		err = dec.Decode(&result)
		if err != nil {
			return user, errors.Wrapf(err, "errored decoding JSON from requests")
		}
		// the reason we add all of these underscores after this is to prevent
		// go from panic ing in case it fails to convert the value to a string.
		// If it does fail, the value will just be an empty string (which is
		// the best we can do in any case)
		user.Username, _ = result["login"].(string)
		user.Name, _ = result["name"].(string)
		user.URL, _ = result["blog"].(string)
		user.Avatar, _ = result["avatar_url"].(string)
		user.Email, _ = result["email"].(string)
		user.Location, _ = result["location"].(string)
		user.Bio, _ = result["bio"].(string)
		return user, nil
	}
	return user, errors.Errorf("unknown service %q to ask user informations from", service)
}

func saveUserInformation(token string, service string, user User) (int, error) {
	sql := `
	INSERT INTO users (token, username, avatar, name, bio, url, email, location, service)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (email, service) DO
	UPDATE SET username=$2, avatar=$3, name=$4, bio=$5, url=$6, email=$7, location=$8
	RETURNING (id)
	`
	var id int
	log.Printf("Save user to database: %v", user)
	err := db.DB().QueryRow(sql, token, user.Username, user.Avatar, user.Name,
		user.Bio, user.URL, user.Email, user.Location, service).Scan(&id)
	return id, errors.Wrapf(err, "errored executing request")
}