package radix_test

import (
	"net/http"
	"testing"

	"github.com/elmq0022/krillin/internal/radix"
	"github.com/elmq0022/krillin/router"
)

func TestNewRadix(t *testing.T) {

	path := "/url/path/to/resource"
	method := http.MethodGet
	handler := 1

	routes := []router.Route[int]{
		{Path: path, Method: method, Handler: handler},
		{Path: "/foo/bar/baz", Method: http.MethodPatch, Handler: 2},
	}

	r, _ := radix.New(routes)

	if r.ChildNPrefix(0) != path {
		t.Fatalf("want: %s, got %s", path, r.ChildNPrefix(0))
	}

	if r.ChildNPrefix(1) != "/foo/bar/baz" {
		t.Fatalf("want: %s, got %s", "/foo/bar/baz", r.ChildNPrefix(1))
	}
}
