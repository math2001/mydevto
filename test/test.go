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
func MakeRequest(method, target string, body io.Reader,
	handler http.HandlerFunc, statuscode int) (*httptest.ResponseRecorder, error) {

	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != statuscode {
		return nil, errors.Errorf("Wrong status code: got %d, want %d\n%q", rr.Code,
			statuscode, readbody(rr.Body))
	}
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		return nil, errors.Errorf("Wrong Content-Type header:got %s, want %s\n%s",
			"application/json", ctype, readbody(rr.Body))
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
