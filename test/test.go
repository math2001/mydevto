package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

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
