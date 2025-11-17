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
	got, _ := r.Lookup(method, path)
	if got != handler {
		t.Fatalf("want %d, got %d", handler, got)
	}

	got2, _ := r.Lookup(http.MethodPatch, "/foo/bar/baz")
	if got2 != 2 {
		t.Fatalf("want %d, got %d", 2, got2)
	}
}
