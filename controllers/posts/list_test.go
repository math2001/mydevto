package posts

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test"
	"github.com/math2001/mydevto/test/testdb"
)

func TestListNoFilter(t *testing.T) {
	rr, err := test.MakeRequest("GET", "/api/posts/list", nil, list)
	if err != nil {
		t.Fatal(err)
	}
	var actual []db.Post
	var text bytes.Buffer
	tee := io.TeeReader(rr.Body, &text)
	dec := json.NewDecoder(tee)
	err = dec.Decode(&actual)
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

func TestListLimit(t *testing.T) {
	// rr, msg := test.MakeRequest("GET", "/api/posts/list", nil, list)

}
