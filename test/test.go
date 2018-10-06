package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gbrlsnchs/jwt"
	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test/testdb"
	"github.com/pkg/errors"
)

// cheating the system: we pretend the user is already logged in, and this is
// the token he has stored as cookie
var token string

func init() {
	key := os.Getenv("JWTSECRET")
	if key == "" {
		log.Fatal("$JWTSECRET needs to be specified")
	}
	hs256 := jwt.NewHS256(key)
	jot := &db.JWTToken{
		JWT:  &jwt.JWT{},
		User: &testdb.Users[0], // we pretend to be the user number 1
	}
	jot.SetAlgorithm(hs256)
	jot.SetKeyID("kid")
	payload, err := jwt.Marshal(jot)
	if err != nil {
		log.Fatalf("could not marshal jot: %s", err)
	}
	b, err := hs256.Sign(payload)
	if err != nil {
		log.Fatalf("could not sign payload: %s", err)
	}
	token = string(b)
}

// returns the body of the response, or the error, as a string
// It is used to present errors when test fail
func readbody(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Sprintf("Couldn't read body: %s", err)
	}
	return string(b)
}

// MakeRequest makes a request to a handler and returns the recorded response
func MakeRequest(server *httptest.Server, method, target string, body io.Reader,
	statuscode int) (*http.Response, error) {

	req, err := http.NewRequest(method, server.URL+target, body)
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(&http.Cookie{
		Name:  api.JWT,
		Value: token,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not make request")
	}
	rr, err := server.Client().Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request")
	}

	if status := rr.StatusCode; status != statuscode {
		return nil, errors.Errorf("Wrong status code: got %d, want %d\n%q",
			rr.StatusCode, statuscode, readbody(rr.Body))
	}
	if ctype := rr.Header.Get("Content-Type"); ctype != "application/json" {
		return nil, errors.Errorf("Wrong Content-Type header: got %q, want %q\n%s",
			ctype, "application/json", readbody(rr.Body))
	}
	return rr, nil
}

// Decode decodes a JSON response into to and returns any error that occurs
func Decode(r io.Reader, to interface{}) error {
	var text bytes.Buffer
	tee := io.TeeReader(r, &text)
	dec := json.NewDecoder(tee)
	err := dec.Decode(&to)
	if err != nil {
		body, err := ioutil.ReadAll(&text)
		if err != nil {
			return errors.Wrapf(err, "Couldn't read duplicated response body")
		}
		return errors.Wrapf(err, "Couldn't decode response body in %q", body)
	}
	return nil
}
