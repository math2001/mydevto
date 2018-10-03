package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

// MakeRequest makes a request to a handler and returns the recorded response
func MakeRequest(method, target string, body io.Reader,
	handler http.HandlerFunc) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		return nil, errors.Errorf("Wrong status code:\ngot %d,\nwant %d", rr.Code,
			http.StatusOK)
	}
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		return nil, errors.Errorf("Wrong Content-Type header:\ngot %s,\nwant %s",
			"application/json", ctype)
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
			return errors.Wrapf(err, "Couldn't readel duplicated response body")
		}
		return errors.Wrapf(err, "Couldn't decode response body in %q", body)
	}
	return nil
}
