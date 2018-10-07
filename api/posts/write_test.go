package posts_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/test"
	"github.com/math2001/mydevto/test/testdb"
)

// this is bad... Tests aren't indenpendent...

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
	_, err = db.ExecContext(context.Background(),
		"DELETE FROM posts WHERE id=$1", msg.ID)
	if err != nil {
		t.Fatalf("could not clean up post: %s", err)
	}
}

func TestWriteUpdate(t *testing.T) {
	updated := time.Now()
	update := db.Post{
		Title:   "test update",
		Content: "updating a post...",
	}
	post := testdb.Posts[0]
	id := post.ID
	vals := url.Values{}
	vals.Add("id", strconv.Itoa(id))
	vals.Add("title", update.Title)
	vals.Add("content", update.Content)
	body := strings.NewReader(vals.Encode())
	res, err := test.MakeRequest(server, "POST", "/api/posts/write", body,
		http.StatusOK)
	if err != nil {
		t.Fatalf("could not make POST request: %s", err)
	}
	var msg struct {
		Type    string
		Message string
	}
	if err = test.Decode(res.Body, &msg); err != nil {
		log.Fatalf("could not decode response: %s", err)
	}
	if msg.Type != "success" || msg.Message == "post successfully updated" {
		t.Fatalf("invalid response: %#v", msg)
	}

	// make the get request to see if the post has been updated
	res, err = test.MakeRequest(server, "GET",
		fmt.Sprintf("/api/posts/get?id=%d", id), nil, http.StatusOK)
	if err != nil {
		t.Fatalf("could not make GET request: %s", err)
	}
	var result db.Post
	if err = test.Decode(res.Body, &result); err != nil {
		t.Fatalf("could not decode post (id %d): %s", id, err)
	}
	if !update.Equals(result) {
		t.Fatalf("update and response don't match:\nwant: %#v\ngot:  %#q",
			update, result)
	}
	if !updated.Before(result.Updated) {
		t.Fatalf("post updated field invalid: %s should be before %s", updated,
			result.Updated)
	}

	// clean up the post
	_, err = db.ExecContext(context.Background(), `UPDATE posts SET
	userid=$1, title=$2, content=$3, updated=$4, written=$5
	WHERE id=$6`, post.User.ID, post.Title, post.Content, post.Updated,
		post.Written, id)
	if err != nil {
		log.Fatalf("could not clean up post: %s", err)
	}
}

func TestWriteErrors(t *testing.T) {
	var params = []url.Values{
		{},
		{"title": {"valid"}, "conte": {"no content"}},
		{"title": {"valid"}, "content": {"no content"}, "id": {"asd"}},
	}
	for _, args := range params {
		args := args
		t.Run(args.Encode(), func(t *testing.T) {
			body := strings.NewReader(args.Encode())
			res, err := test.MakeRequest(server, "POST", "/api/posts/write", body,
				http.StatusBadRequest)
			if err != nil {
				t.Fatalf("could not make request for args %q: %s", args.Encode(), err)
			}
			var text map[string]string
			if err = test.Decode(res.Body, &text); err != nil {
				t.Fatalf("could not decode response: %s", err)
			}
			if !reflect.DeepEqual(text, map[string]string{
				"type":    "error",
				"message": "Invalid request form data",
			}) {
				t.Fatalf("expected invalid data message, got: %v", text)
			}
		})
	}

}
