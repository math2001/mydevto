package posts_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test"
)

func TestInsert(t *testing.T) {
	insert := time.Now()
	post := db.Post{
		Title:   "test insert post",
		Content: "This is just testing whether inserting a post works",
	}
	vals := url.Values{}
	vals.Add("title", post.Title)
	vals.Add("content", post.Content)
	body := strings.NewReader(vals.Encode())
	res, err := test.MakeRequest(server, "POST", "/api/posts/write", body,
		http.StatusOK)
	if err != nil {
		t.Fatalf("could not make POST request: %s", err)
	}
	var msg struct {
		Type    string
		Message string
		ID      int
	}
	if err = test.Decode(res.Body, &msg); err != nil {
		t.Fatalf("could not decode response: %s", err)
	}
	if msg.Type != "success" || msg.Message != "post successfully inserted" {
		t.Fatalf("invalid response: %#v", msg)
	}
	// Now, we make a get request asking for the post, and then checking if it
	// matches
	res, err = test.MakeRequest(server, "GET",
		fmt.Sprintf("/api/posts/get?id=%d", msg.ID), nil, http.StatusOK)
	if err != nil {
		t.Fatalf("could not make GET request: %s", err)
	}
	var result db.Post
	if err = test.Decode(res.Body, &result); err != nil {
		log.Fatalf("could not decode post (id %d): %s", msg.ID, err)
	}
	if !post.Equals(result) {
		t.Fatalf("post does not match\nwant: %#v\ngot:  %#v", post, result)
	}
	if !insert.Before(result.Written) {
		t.Fatalf("post written field invalid: %s should be before %s", insert,
			result.Written)
	}

	// clean up the post
	_, err = db.DB().Exec("DELETE FROM posts WHERE id=$1", msg.ID)
	if err != nil {
		t.Fatalf("could not clean up post: %s", err)
	}
}
