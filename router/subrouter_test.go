package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
)

func TestSubRouter(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("%v", err)
	}

	api := r.Group("/api/v1/")
	wantStatus := 200
	wantBody := "bar"

	api.GET("/foo", NewTestHandler(wantStatus, wantBody))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/foo", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != wantStatus {
		t.Fatalf("want %d, got %d", wantStatus, rr.Code)
	}

	if rr.Body.String() != wantBody {
		t.Fatalf("want %s, got %s", wantBody, rr.Body.String())
	}

	// if spy.Err != wantErr {
	// 	t.Fatalf("want %v, got %v", &wantErr, spy.Err)
	// }
}
