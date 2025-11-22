package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

func TestWithMiddleware(t *testing.T) {
	// r, _ := router.New(NewSpyAdapter(&SpyAdapterRecord{}), router.WithMiddleware())

}

func TestWithNotFound(t *testing.T) {
	testNotFound := func(r *http.Request) (types.Response, error) {
		return types.Response{
				Status: http.StatusNotFound,
				Body:   "test not found"},
			nil
	}
	spy := SpyAdapterRecord{}
	r, _ := router.New(NewSpyAdapter(&spy), router.WithNotFound(testNotFound))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rec, req)

	if spy.Status != http.StatusNotFound {
		t.Fatalf("want %d got %d", http.StatusNotFound, spy.Status)
	}

	if spy.Body != "test not found" {
		t.Fatalf("want %s, got %s", "test not found", spy.Body)
	}
}
