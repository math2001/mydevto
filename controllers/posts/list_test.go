package posts

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListNoFilter(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/posts/list", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(list)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: want %d, got %d", http.StatusOK, rr.Code)
	}
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Wrong Content-Type header: want 'application/json' %q",
			ctype)
	}
}
