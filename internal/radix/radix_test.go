package radix_test

import (
	"net/http"
	"testing"

	"github.com/elmq0022/kami/internal/radix"
	"github.com/elmq0022/kami/types"
)

func MakeTestHandler(value any) types.Handler {
	return func(req *http.Request) (int, any, error) {
		return 0, value, nil
	}
}

func ReadTestHandler(h types.Handler) any {
	fakeReq, _ := http.NewRequest(http.MethodGet, "", nil)
	_, value, _ := h(fakeReq)
	return value
}

func TestNewRadix(t *testing.T) {

	path := "/foo/bar/baz"
	method := http.MethodGet
	handler := MakeTestHandler(1)

	routes := types.Routes{
		{Path: path, Method: method, Handler: handler},
		{Path: "/foo/bar/baz2", Method: http.MethodGet, Handler: MakeTestHandler(2)},
		{Path: "/foo/bar/:id", Method: http.MethodGet, Handler: MakeTestHandler(3)},
	}

	r, _ := radix.New(routes)

	h, _, _ := r.Lookup(method, path)
	if got := ReadTestHandler(h); got != 1 {
		t.Fatalf("want %d, got %d", 1, got)
	}

	h, _, _ = r.Lookup(http.MethodGet, "/foo/bar/baz2")
	if got := ReadTestHandler(h); got != 2 {
		t.Fatalf("want %d, got %d", 2, got)
	}

	h, params, _ := r.Lookup(http.MethodGet, "/foo/bar/42")
	if got := ReadTestHandler(h); got != 3 {
		t.Fatalf("want %d, got %d", 3, got)
	}

	param, ok := params["id"]
	if !ok {
		t.Fatal("could not retrive the paramter")
	}
	if param != "42" {
		t.Fatalf("want 42, got %s", param)
	}
}
