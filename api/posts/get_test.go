package posts_test

import (
	"reflect"
	"testing"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test"
	"github.com/math2001/mydevto/test/testdb"
)

func TestGetValid(t *testing.T) {
	var tt = []struct {
		name       string
		statuscode int
		uri        string
		body       db.Post
	}{
		{"id=2", 200, "/api/posts/get?id=2", testdb.Posts[1]},
		{"id=1", 200, "/api/posts/get?id=1", testdb.Posts[0]},
	}
	for _, tc := range tt {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			rr, err := test.MakeRequest(server, "GET", tc.uri, nil, tc.statuscode)
			if err != nil {
				t.Fatal(err)
			}
			var post db.Post
			if err = test.Decode(rr.Body, &post); err != nil {
				t.Fatal(err)
			}
			if !post.Equals(tc.body) {
				t.Errorf("Posts don't match: \ngot:  %#v\nwant: %#v", post,
					tc.body)
			}
		})
	}
}

func TestGetError(t *testing.T) {
	var tt = []struct {
		name       string
		statuscode int
		uri        string
		body       map[string]string
	}{
		{"non existing id", 400, "/api/posts/get?id=6", map[string]string{
			"type":    "error",
			"message": "No post found with id 6",
		}},
		{"negative id", 400, "/api/posts/get?id=-12", map[string]string{
			"type":    "error",
			"message": "No post found with id -12",
		}},
		{"float", 400, "/api/posts/get?id=12.3", map[string]string{
			"type":    "error",
			"message": "Couldn't convert id \"12.3\" to integer",
		}},
		{"text", 400, "/api/posts/get?id=hello", map[string]string{
			"type":    "error",
			"message": "Couldn't convert id \"hello\" to integer",
		}},
		{"text", 400, "/api/posts/get", map[string]string{
			"type":    "error",
			"message": "Invalid id. Empty strings not allowed",
		}},
	}
	for _, tc := range tt {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			rr, err := test.MakeRequest(server, "GET", tc.uri, nil, tc.statuscode)
			if err != nil {
				t.Fatal(err)
			}
			var result map[string]string
			if err = test.Decode(rr.Body, &result); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, tc.body) {
				t.Errorf("Posts don't match: \ngot:  %#v\nwant: %#v", result,
					tc.body)
			}
		})
	}
}
