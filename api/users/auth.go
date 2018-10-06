package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"
	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/services/uli"
	"github.com/pkg/errors"
)

// auth is called by the service with 'service' as the name in the
// parameters, and finishes the auth flow. It's the 'callback'
func auth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	service := query.Get("service")
	if service != "github" {
		uli.Printf(r, "Invalid service %q trying to authenticate @ usersAuth", service)
		api.Error(w, r, http.StatusBadRequest, "Invalid service")
		return
	}
	sessioncode := query.Get("code")
	if sessioncode == "" {
		uli.Printf(r, "No session code in URL parameter @ usersAuth#%q", service)
		api.InternalError(w, r)
		return
	}
	servicetoken, err := getToken(sessioncode)
	if err != nil {
		uli.Printf(r, "Errored getting token: %s", err)
		api.InternalError(w, r)
		return
	}
	user, err := retrieveUserInformation(servicetoken, api.ServiceGithub)
	if err != nil {
		uli.Printf(r, "Errored retrieving user information from token: %s", err)
		api.InternalError(w, r)
		return
	}
	id, err := saveUserInformation(servicetoken, api.ServiceGithub, user)
	if err != nil {
		uli.Printf(r, "Errored saving user information to database: %s", err)
		api.InternalError(w, r)
		return
	}
	user.ID = id
	jot := &db.JWTToken{
		JWT:  &jwt.JWT{},
		User: &user,
	}
	jot.SetAlgorithm(jwtsigner)
	jot.SetKeyID("kid")
	payload, err := jwt.Marshal(jot)
	if err != nil {
		uli.Printf(r, "could not marshal JSON token: %s", err)
		api.InternalError(w, r)
		return
	}
	token, err := jwtsigner.Sign(payload)
	if err != nil {
		uli.Printf(r, "could not sign payload: %s", err)
		api.InternalError(w, r)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     api.JWT,
		Value:    string(token),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour), // a month
	})
	fmt.Fprintf(w, "<script>window.close()</script>")
}

// getToken retrieves the token from a temporary code (this is part of the
// oauth flow)
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
	res, err := api.HTTPClient.Do(req)
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

// Gets the user information from a token
func retrieveUserInformation(token string, service string) (db.User, error) {
	var user db.User
	if service == api.ServiceGithub {
		req, err := http.NewRequest("GET", "https://api.github.com/user?access_token="+token, nil)
		if err != nil {
			return user, errors.Wrapf(err, "errored creating request")
		}
		req.Header.Add("Accept", "application/json")
		res, err := api.HTTPClient.Do(req)
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

// Saves the user information into the database, returning the ID of that user,
// and an error, if any.
func saveUserInformation(token string, service string, user db.User) (int, error) {
	sql := `
	INSERT INTO users (token, username, avatar, name, bio, url, email, location, service)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (email, service) DO
	UPDATE SET username=$2, avatar=$3, name=$4, bio=$5, url=$6, email=$7, location=$8
	RETURNING (id)
	`
	var id int
	err := db.DB().QueryRow(sql, token, user.Username, user.Avatar, user.Name,
		user.Bio, user.URL, user.Email, user.Location, service).Scan(&id)
	return id, errors.Wrapf(err, "errored executing request")
}
