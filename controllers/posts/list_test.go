package posts

import (
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
	if err = test.Decode(rr.Body, &actual); err != nil {
		t.Fatal(err)
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
	rr, err := test.MakeRequest("GET", "/api/posts/list?limit=1", nil, list)
	if err != nil {
		t.Fatal(err)
	}
	var actual []db.Post
	if err = test.Decode(rr.Body, &actual); err != nil {
		t.Fatal(err)
	}
	if len(actual) != 1 {
		t.Errorf("Response length didn't match: got %d, want %d",
			len(actual), 1)
	}
	if !actual[0].Equals(testdb.Posts[0]) {
		t.Errorf("Post didn't match:\n%v\n%v", actual[0], testdb.Posts[0])
	}
}
