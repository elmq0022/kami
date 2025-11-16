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
	}

	r, _ := radix.New(routes)

	if r.FirstPrefix() != path {
		t.Fatalf("want: %s, got %s", path, r.FirstPrefix())
	}
}
