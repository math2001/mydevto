package test

import (
	"io"
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
