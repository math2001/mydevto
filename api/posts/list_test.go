package posts_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/math2001/mydevto/router"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test"
	"github.com/math2001/mydevto/test/testdb"
)

var server *httptest.Server

func TestMain(t *testing.M) {
	server = httptest.NewServer(router.Router())
	code := t.Run()
	server.Close()
	os.Exit(code)
}

func TestListNoFilter(t *testing.T) {
	rr, err := test.MakeRequest(server, "GET", "/api/posts/list", nil,
		http.StatusOK)
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
	var found bool
	for _, post := range actual {
		found = false
		for _, p := range testdb.Posts {
			if p.ID == post.ID {
				found = true
				if !post.Equals(p) {
					t.Errorf("Post didn't match: \n%v\n%v", post, p)
				}
			}
		}
		if found == false {
			t.Fatalf("could not find post id %d in result", post.ID)
		}
	}
}

func TestListLimit(t *testing.T) {
	rr, err := test.MakeRequest(server, "GET", "/api/posts/list?limit=1", nil,
		http.StatusOK)
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
	for _, post := range testdb.Posts {
		if post.ID == actual[0].ID {
			if !actual[0].Equals(post) {
				t.Errorf("Post didn't match: \n%v\n%v", actual[0], post)
			}
			return
		}
	}
}

func TestListUserid(t *testing.T) {
	rr, err := test.MakeRequest(server, "GET", "/api/posts/list?userid=1", nil,
		http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	var actual []db.Post
	if err = test.Decode(rr.Body, &actual); err != nil {
		t.Fatal(err)
	}
	if len(actual) != 2 {
		t.Errorf("Response length didn't match: got %d, want %d",
			len(actual), 2)
	}
	var found bool
	for _, post := range actual {
		found = false
		for _, p := range testdb.Posts {
			if p.ID == post.ID {
				found = true
				if !post.Equals(p) {
					t.Errorf("Post didn't match: \n%v\n%v", post, p)
				}
			}
		}
		if found == false {
			t.Fatalf("could not find post id %d in result", post.ID)
		}
	}
}

func TestListLimitUserid(t *testing.T) {
	rr, err := test.MakeRequest(server, "GET", "/api/posts/list?userid=1&limit=1",
		nil, http.StatusOK)
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
	for _, post := range testdb.Posts {
		if post.ID == actual[0].ID {
			if !actual[0].Equals(post) {
				t.Errorf("Post didn't match: \n%v\n%v", actual[0], post)
			}
			return
		}
	}
}
