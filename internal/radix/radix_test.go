package radix_test

import (
	"net/http"
	"testing"

	"github.com/elmq0022/kami/internal/radix"
	"github.com/elmq0022/kami/types"
)

type testRenerable struct {
	Value any
}

func (tr *testRenerable) Render(w http.ResponseWriter) {
}

func MakeTestHandler(value any) types.Handler {
	return func(req *http.Request) types.Renderable {
		return &testRenerable{Value: value}
	}
}

func ReadTestHandler(h types.Handler) any {
	fakeReq, _ := http.NewRequest(http.MethodGet, "", nil)
	resp := h(fakeReq).(*testRenerable)
	return resp.Value
}

func TestRadix_AddRoute_Validation(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			name:      "empty parameter name",
			path:      "/user/:/posts",
			wantError: true,
		},
		{
			name:      "single char parameter name",
			path:      "/user/:a/posts",
			wantError: false,
		},
		{
			name:      "empty wildcard name",
			path:      "/static/*",
			wantError: true,
		},
		{
			name:      "single char wildcard name",
			path:      "/static/*p",
			wantError: false,
		},
		{
			name:      "valid two char parameter",
			path:      "/user/:id/posts",
			wantError: false,
		},
		{
			name:      "valid two char wildcard",
			path:      "/static/*fp",
			wantError: false,
		},
		{
			name:      "path without leading slash",
			path:      "user/profile",
			wantError: true,
		},
		{
			name:      "empty path",
			path:      "",
			wantError: true,
		},
		{
			name:      "wildcard in middle position",
			path:      "/static/*/more",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := radix.New()
			err := r.AddRoute(http.MethodGet, tt.path, MakeTestHandler("test"))

			if tt.wantError && err == nil {
				t.Fatalf("expected error for path %q, got nil", tt.path)
			}
			if !tt.wantError && err != nil {
				t.Fatalf("expected no error for path %q, got %v", tt.path, err)
			}
		})
	}
}

func TestRadix_DuplicateParameterNames(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "duplicate parameter names in path",
			path:      "/user/:id/posts/:id",
			wantError: true,
		},
		{
			name:      "duplicate parameter names with different structure",
			path:      "/api/:version/users/:id/posts/:version",
			wantError: true,
		},
		{
			name:      "unique parameter names",
			path:      "/user/:userId/posts/:postId",
			wantError: false,
		},
		{
			name:      "three parameters all unique",
			path:      "/api/:version/users/:userId/posts/:postId",
			wantError: false,
		},
		{
			name:      "duplicate wildcard and param with same name",
			path:      "/user/:path/files/*path",
			wantError: true,
		},
		{
			name:      "single parameter no duplicates",
			path:      "/user/:id",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := radix.New()
			err := r.AddRoute(http.MethodGet, tt.path, MakeTestHandler("test"))

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error for path %q, got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for path %q, got %v", tt.path, err)
				}
			}
		})
	}
}

func TestRadix_Lookup(t *testing.T) {
	tests := []struct {
		name       string
		routes     types.Routes
		method     string
		path       string
		wantValue  any
		wantParams map[string]string
		wantFound  bool
	}{
		// Static routes
		{
			name: "static root",
			routes: types.Routes{
				{Path: "/", Method: http.MethodGet, Handler: MakeTestHandler("root")},
			},
			method:    http.MethodGet,
			path:      "/",
			wantValue: "root",
			wantFound: true,
		},
		{
			name: "static nested",
			routes: types.Routes{
				{Path: "/foo/bar", Method: http.MethodGet, Handler: MakeTestHandler("bar")},
			},
			method:    http.MethodGet,
			path:      "/foo/bar",
			wantValue: "bar",
			wantFound: true,
		},

		// Param routes
		{
			name: "single param",
			routes: types.Routes{
				{Path: "/user/:id", Method: http.MethodGet, Handler: MakeTestHandler("user")},
			},
			method:     http.MethodGet,
			path:       "/user/alice",
			wantValue:  "user",
			wantParams: map[string]string{"id": "alice"},
			wantFound:  true,
		},
		{
			name: "two params",
			routes: types.Routes{
				{Path: "/user/:uid/post/:pid", Method: http.MethodGet, Handler: MakeTestHandler("post")},
			},
			method:     http.MethodGet,
			path:       "/user/alice/post/42",
			wantValue:  "post",
			wantParams: map[string]string{"uid": "alice", "pid": "42"},
			wantFound:  true,
		},

		// Wildcards
		{
			name: "wildcard match",
			routes: types.Routes{
				{Path: "/static/*path", Method: http.MethodGet, Handler: MakeTestHandler("static")},
			},
			method:     http.MethodGet,
			path:       "/static/js/app.js",
			wantValue:  "static",
			wantParams: map[string]string{"path": "js/app.js"},
			wantFound:  true,
		},

		// Conflicting routes
		{
			name: "param vs static conflict",
			routes: types.Routes{
				{Path: "/user/list", Method: http.MethodGet, Handler: MakeTestHandler("list")},
				{Path: "/user/:id", Method: http.MethodGet, Handler: MakeTestHandler("detail")},
			},
			method:     http.MethodGet,
			path:       "/user/list",
			wantValue:  "list", // static should win
			wantParams: map[string]string{},
			wantFound:  true,
		},
		{
			name: "param wins if no static",
			routes: types.Routes{
				{Path: "/user/:id", Method: http.MethodGet, Handler: MakeTestHandler("detail")},
			},
			method:     http.MethodGet,
			path:       "/user/bob",
			wantValue:  "detail",
			wantParams: map[string]string{"id": "bob"},
			wantFound:  true,
		},

		// Method mismatch
		{
			name: "wrong method",
			routes: types.Routes{
				{Path: "/foo", Method: http.MethodGet, Handler: MakeTestHandler("get")},
			},
			method:    http.MethodPost,
			path:      "/foo",
			wantFound: false,
		},

		// Not found
		{
			name: "missing route",
			routes: types.Routes{
				{Path: "/foo/bar", Method: http.MethodGet, Handler: MakeTestHandler("bar")},
			},
			method:    http.MethodGet,
			path:      "/foo/baz",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := radix.New()
			m := tt.routes[0].Method
			p := tt.routes[0].Path
			h := tt.routes[0].Handler
			r.AddRoute(m, p, h)
			if err != nil {
				t.Fatalf("failed to create radix: %v", err)
			}

			h, params, found := r.Lookup(tt.method, tt.path)
			if found != tt.wantFound {
				t.Fatalf("expected found=%v, got %v", tt.wantFound, found)
			}
			if !found {
				return
			}

			got := ReadTestHandler(h)
			if got != tt.wantValue {
				t.Fatalf("expected value %v, got %v", tt.wantValue, got)
			}

			for k, v := range tt.wantParams {
				if params[k] != v {
					t.Fatalf("expected param %s=%s, got %s", k, v, params[k])
				}
			}
		})
	}
}
