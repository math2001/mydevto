package posts

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test/testdb"
)

func TestListNoFilter(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/posts/list", nil)
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
	var actual []db.Post
	var text bytes.Buffer
	tee := io.TeeReader(rr.Body, &text)
	dec := json.NewDecoder(tee)
	err := dec.Decode(&actual)
	if err != nil {
		t.Errorf("Couldn't decode response body: %s", err)
		b, err := ioutil.ReadAll(&text)
		if err != nil {
			t.Errorf("Couldn't read from duplicated body: %s", err)
		}
		t.Logf("Body: %q", string(b))
		t.Fatal()
	}
	if len(testdb.Posts) != len(actual) {
		t.Fatalf("Response length didn't match: want %d, got %d",
			len(testdb.Posts), len(actual))
	}
	for i, post := range actual {
		if !post.Equals(testdb.Posts[i]) {
			t.Errorf("Post didn't match: \n%v\n%v", post, testdb.Posts[i])
		}
	}
}
