package router_test

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

func NewTestHandler(status int, body string) types.Handler {
	return func(req *http.Request) types.Responder {
		return &testResponder{Status: status, Body: body}
	}
}

func TestRouter_RoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		wantErr    error
		wantParams map[string]string
		callPath   string
	}{
		// Static routes
		{name: "root", method: http.MethodGet, path: "/", wantStatus: http.StatusOK, wantBody: "root", wantErr: nil, wantParams: map[string]string{}, callPath: "/"},
		{name: "about", method: http.MethodGet, path: "/about", wantStatus: http.StatusOK, wantBody: "about", wantErr: nil, wantParams: map[string]string{}, callPath: "/about"},

		// Param routes
		{name: "book by id", method: http.MethodGet, path: "/book/:id", wantStatus: http.StatusOK, wantBody: "books", wantErr: nil, wantParams: map[string]string{"id": "lifeOfPi"}, callPath: "/book/lifeOfPi"},
		{name: "user post", method: http.MethodGet, path: "/user/:userId/post/:postId", wantStatus: http.StatusOK, wantBody: "post", wantErr: nil, wantParams: map[string]string{"userId": "alice", "postId": "42"}, callPath: "/user/alice/post/42"},

		// Overlapping param routes
		{name: "user list", method: http.MethodGet, path: "/user/list", wantStatus: http.StatusOK, wantBody: "user list", wantErr: nil, wantParams: map[string]string{}, callPath: "/user/list"},
		{name: "user detail", method: http.MethodGet, path: "/user/:id", wantStatus: http.StatusOK, wantBody: "user detail", wantErr: nil, wantParams: map[string]string{"id": "bob"}, callPath: "/user/bob"},

		// Wildcard routes
		{name: "static js", method: http.MethodGet, path: "/static/*path", wantStatus: http.StatusOK, wantBody: "static", wantErr: nil, wantParams: map[string]string{"path": "js/app.js"}, callPath: "/static/js/app.js"},
		{name: "static css", method: http.MethodGet, path: "/static/*path", wantStatus: http.StatusOK, wantBody: "static", wantErr: nil, wantParams: map[string]string{"path": "css/main.css"}, callPath: "/static/css/main.css"},

		// Method mismatch (should not match)
		{name: "wrong method", method: http.MethodPost, path: "/about", wantStatus: http.StatusNotFound, wantBody: "Not Found", wantErr: nil, wantParams: map[string]string{}, callPath: "/about"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, err := router.New()
			if err != nil {
				t.Fatalf("failed to create router: %v", err)
			}

			var gotParams map[string]string
			handler := func(req *http.Request) types.Responder {
				gotParams = router.GetParams(req.Context())
				return NewTestHandler(tt.wantStatus, tt.wantBody)(req)
			}

			r.GET(tt.path, handler)

			req := httptest.NewRequest(tt.method, tt.callPath, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if tt.wantBody != rr.Body.String() {
				t.Fatalf("body: want %v, got %v", tt.wantBody, rr.Body.String())
			}

			if tt.wantStatus != rr.Code {
				t.Fatalf("status: want %d, got %d", tt.wantStatus, rr.Code)
			}

			if !maps.Equal(tt.wantParams, gotParams) {
				t.Fatalf("params: want %v got %v", tt.wantParams, gotParams)
			}
		})
	}
}

func TestRouter_CannotAddRoutesAfterStarted(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Add a route before starting
	r.GET("/before", NewTestHandler(http.StatusOK, "before"))

	// Simulate the router being started by making a request
	req := httptest.NewRequest(http.MethodGet, "/before", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Now try to add a route after the router has handled a request
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatal("expected panic when adding route after router started, got nil")
		} else {
			panicMsg := rec.(string)
			expectedMsg := "cannot register path: /after since the router is running"
			if panicMsg != expectedMsg {
				t.Fatalf("unexpected panic message: got %q, want %q", panicMsg, expectedMsg)
			}
		}
	}()

	r.GET("/after", NewTestHandler(http.StatusOK, "after"))
}
