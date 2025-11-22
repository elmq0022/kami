package router_test

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

type SpyAdapterRecord struct {
	Status int
	Body   any
	Err    error
	Params map[string]string
}

func NewSpyAdapter(record *SpyAdapterRecord) types.Adapter {
	return func(w http.ResponseWriter, r *http.Request, h types.Handler) {
		resp, err := h(r)
		status := resp.Status
		body := resp.Body
		record.Status = status
		record.Body = body
		record.Err = err
		record.Params = router.GetParams(r.Context())
	}
}

func NewTestHandler(status int, body any, err error) types.Handler {
	return func(req *http.Request) (types.Response, error) {
		return types.Response{Status: status, Body: body}, err
	}
}

func TestRouter_RoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   any
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
			record := SpyAdapterRecord{}
			r, err := router.New(NewSpyAdapter(&record))
			if err != nil {
				t.Fatalf("failed to create router: %v", err)
			}

			r.GET(tt.path, NewTestHandler(tt.wantStatus, tt.wantBody, tt.wantErr))

			req := httptest.NewRequest(tt.method, tt.callPath, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if tt.wantBody != record.Body {
				t.Fatalf("body: want %v, got %v", tt.wantBody, record.Body)
			}

			if tt.wantStatus != record.Status {
				t.Fatalf("status: want %d, got %d", tt.wantStatus, record.Status)
			}

			if tt.wantErr != record.Err {
				t.Fatalf("error: want %v got %v", tt.wantErr, record.Err)
			}

			if !maps.Equal(tt.wantParams, record.Params) {
				t.Fatalf("params: want %v got %v", tt.wantParams, record.Params)
			}
		})
	}
}
