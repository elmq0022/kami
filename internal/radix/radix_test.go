package radix_test

import (
	"net/http"
	"testing"

	"github.com/elmq0022/krillin/internal/radix"
	"github.com/elmq0022/krillin/router"
)

func TestNewRadix(t *testing.T) {

	path := "/foo/bar/baz"
	method := http.MethodGet
	handler := func(req *http.Request) (int, any, error) { return 200, 1, nil }

	routes := router.Routes{
		{Path: path, Method: method, Handler: handler},
		{Path: "/foo/bar/baz2", Method: http.MethodPatch, Handler: func(req *http.Request) (int, any, error) { return 200, 2, nil }},
	}

	r, _ := radix.New(routes)
	fakeReq, _ := http.NewRequest(http.MethodGet, "", nil)

	h, _ := r.Lookup(method, path)
	_, got, _ := h(fakeReq)
	if got != 1 {
		t.Fatalf("want %d, got %d", 1, got)
	}

	h, _ = r.Lookup(http.MethodPatch, "/foo/bar/baz2")
	_, got2, _ := h(fakeReq)
	if got2 != 2 {
		t.Fatalf("want %d, got %d", 2, got2)
	}
}
