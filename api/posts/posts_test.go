package posts_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/math2001/mydevto/router"
)

var server *httptest.Server

func TestMain(t *testing.M) {
	server = httptest.NewServer(router.Router())
	code := t.Run()
	server.Close()
	os.Exit(code)
}
