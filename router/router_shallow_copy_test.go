package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

// TestShallowCopy_ReturnsDifferentInstance verifies that shallowCopy returns a new Router instance
func TestShallowCopy_ReturnsDifferentInstance(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Note: shallowCopy is unexported, but we can test it via Use() which calls it
	r2 := r.Use()

	if r == r2 {
		t.Fatal("Use() should return a different Router instance")
	}
}

// TestShallowCopy_SharesRadixTree verifies that shallow copies share the same radix tree
func TestShallowCopy_SharesRadixTree(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Register a route on the original router
	handler := func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "shared"}
	}
	r.Prefix("/shared").GET(handler)

	// Create a shallow copy
	r2 := r.Use()

	// The route registered on r should be accessible via r2 (shared radix tree)
	req := httptest.NewRequest(http.MethodGet, "/shared", nil)
	rr := httptest.NewRecorder()
	r2.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "shared" {
		t.Fatalf("expected body 'shared', got %q", rr.Body.String())
	}
}

// TestMiddlewareChaining_ParentNotAffected verifies that child routers don't modify parent middleware
func TestMiddlewareChaining_ParentNotAffected(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Add middleware to parent
	r = r.Use(testMiddleware1)

	// Register a route on parent
	r.Prefix("/parent").GET(testHandler)

	// Create child with additional middleware
	child := r.Use(testMiddleware2)
	child.Prefix("/child").GET(testHandler)

	// Test parent route - should only have middleware1
	t.Run("parent has only its own middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/parent", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		want := "1" // only testMiddleware1
		got := rr.Body.String()
		if got != want {
			t.Errorf("parent route: want %q, got %q", want, got)
		}
	})

	// Test child route - should have both middleware1 and middleware2
	t.Run("child has accumulated middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/child", nil)
		rr := httptest.NewRecorder()
		child.ServeHTTP(rr, req)

		want := "21" // testMiddleware1 + testMiddleware2
		got := rr.Body.String()
		if got != want {
			t.Errorf("child route: want %q, got %q", want, got)
		}
	})
}

// TestMiddlewareChaining_OrderIsCorrect verifies middleware is applied in registration order
func TestMiddlewareChaining_OrderIsCorrect(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Chain middleware using multiple Use() calls
	r1 := r.Use(testMiddleware1)
	r2 := r1.Use(testMiddleware2)
	r3 := r2.Use(testMiddleware3)

	r3.Prefix("/test").GET(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	r3.ServeHTTP(rr, req)

	// Middleware wraps in reverse during add (lines 102-104 in router.go)
	// So mw1 wraps mw2 wraps mw3 wraps handler
	// Result: 3, 2, 1 are appended in that order
	want := "321"
	got := rr.Body.String()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

// TestMiddlewareChaining_MultipleMiddlewareInOneCall verifies multiple middleware in single Use() call
func TestMiddlewareChaining_MultipleMiddlewareInOneCall(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Add multiple middleware at once
	r = r.Use(testMiddleware1, testMiddleware2, testMiddleware3)
	r.Prefix("/test").GET(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	want := "321"
	got := rr.Body.String()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

// TestPrefix_ConcatenatesCorrectly verifies prefix concatenation with proper slash handling
func TestPrefix_ConcatenatesCorrectly(t *testing.T) {
	tests := []struct {
		name           string
		basePrefix     string
		segment        string
		routePath      string
		expectedPath   string
	}{
		{
			name:         "simple concatenation",
			basePrefix:   "",
			segment:      "api",
			routePath:    "/test",
			expectedPath: "/api/test",
		},
		{
			name:         "with leading slash in segment",
			basePrefix:   "",
			segment:      "/api",
			routePath:    "/test",
			expectedPath: "/api/test",
		},
		{
			name:         "with trailing slash in base",
			basePrefix:   "/api/",
			segment:      "v1",
			routePath:    "/test",
			expectedPath: "/api/v1/test",
		},
		{
			name:         "both have slashes",
			basePrefix:   "/api/",
			segment:      "/v1",
			routePath:    "/test",
			expectedPath: "/api/v1/test",
		},
		{
			name:         "nested prefixes",
			basePrefix:   "/api/v1",
			segment:      "users",
			routePath:    "/test",
			expectedPath: "/api/v1/users/test",
		},
		{
			name:         "empty segment returns unchanged",
			basePrefix:   "/api",
			segment:      "",
			routePath:    "/test",
			expectedPath: "/api/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := router.New()
			if err != nil {
				t.Fatalf("failed to create router: %v", err)
			}

			// Apply prefixes
			if tt.basePrefix != "" {
				r = r.Prefix(tt.basePrefix)
			}
			r = r.Prefix(tt.segment)

			handler := func(req *http.Request) types.Responder {
				return &testResponder{Status: http.StatusOK, Body: "ok"}
			}

			// Register with the route path
			r.Prefix(tt.routePath).GET(handler)

			// Try to access at the expected full path
			req := httptest.NewRequest(http.MethodGet, tt.expectedPath, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected route at %q to return 200, got %d", tt.expectedPath, rr.Code)
			}
		})
	}
}

// TestPrefix_EmptySegment verifies empty prefix returns a copy with unchanged prefix
func TestPrefix_EmptySegment(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	r = r.Prefix("/api")
	r2 := r.Prefix("")

	// Should be a different instance
	if r == r2 {
		t.Error("Prefix with empty string should return a new instance")
	}

	// Both should register routes at the same prefix
	r.Prefix("/test1").GET(func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "test1"}
	})

	r2.Prefix("/test2").GET(func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "test2"}
	})

	// Both should be accessible at /api/*
	tests := []struct {
		path string
		want string
	}{
		{"/api/test1", "test1"},
		{"/api/test2", "test2"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rr := httptest.NewRecorder()
		r2.ServeHTTP(rr, req)

		if rr.Body.String() != tt.want {
			t.Errorf("path %q: want %q, got %q", tt.path, tt.want, rr.Body.String())
		}
	}
}

// TestRouterIsolation_SiblingsDontAffectEachOther verifies sibling routers are independent
func TestRouterIsolation_SiblingsDontAffectEachOther(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Create two sibling routers with different middleware
	api := r.Prefix("/api").Use(testMiddleware1)
	admin := r.Prefix("/admin").Use(testMiddleware2)

	api.Prefix("/test").GET(testHandler)
	admin.Prefix("/test").GET(testHandler)

	// Test /api/test - should only have middleware1
	t.Run("api has only its middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		rr := httptest.NewRecorder()
		api.ServeHTTP(rr, req)

		want := "1"
		got := rr.Body.String()
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	// Test /admin/test - should only have middleware2
	t.Run("admin has only its middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		rr := httptest.NewRecorder()
		admin.ServeHTTP(rr, req)

		want := "2"
		got := rr.Body.String()
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})
}

// TestRouterIsolation_ParentNotMutated verifies child modifications don't affect parent
func TestRouterIsolation_ParentNotMutated(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Parent with one middleware
	parent := r.Use(testMiddleware1)
	parent.Prefix("/parent").GET(testHandler)

	// Child adds more middleware
	child := parent.Use(testMiddleware2, testMiddleware3)
	child.Prefix("/child").GET(testHandler)

	// Parent should still only have middleware1
	req := httptest.NewRequest(http.MethodGet, "/parent", nil)
	rr := httptest.NewRecorder()
	parent.ServeHTTP(rr, req)

	want := "1"
	got := rr.Body.String()
	if got != want {
		t.Errorf("parent should not be affected by child changes: want %q, got %q", want, got)
	}
}

// TestComplexChaining_RealWorldScenario simulates the example from the requirements
func TestComplexChaining_RealWorldScenario(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Define middleware that track execution order
	var execOrder []string

	authMW := func(next types.Handler) types.Handler {
		return func(req *http.Request) types.Responder {
			execOrder = append(execOrder, "auth")
			return next(req)
		}
	}

	loggingMW := func(next types.Handler) types.Handler {
		return func(req *http.Request) types.Responder {
			execOrder = append(execOrder, "logging")
			return next(req)
		}
	}

	// Setup the router chain as in the example
	api := r.Prefix("/api").Use(authMW)
	v1 := api.Prefix("/v1")

	v1Handler := func(req *http.Request) types.Responder {
		execOrder = append(execOrder, "v1-handler")
		return &testResponder{Status: http.StatusOK, Body: "v1"}
	}
	v1.GET(v1Handler) // registers GET /api/v1

	users := v1.Prefix("/users").Use(loggingMW)

	listUsersHandler := func(req *http.Request) types.Responder {
		execOrder = append(execOrder, "list-users")
		return &testResponder{Status: http.StatusOK, Body: "users"}
	}
	users.GET(listUsersHandler) // registers GET /api/v1/users

	// Test /api/v1 - should have authMW
	t.Run("v1 endpoint has auth middleware", func(t *testing.T) {
		execOrder = []string{} // reset
		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rr := httptest.NewRecorder()
		v1.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		expectedOrder := []string{"auth", "v1-handler"}
		if len(execOrder) != len(expectedOrder) {
			t.Fatalf("execution order length: want %v, got %v", expectedOrder, execOrder)
		}
		for i, step := range expectedOrder {
			if execOrder[i] != step {
				t.Errorf("step %d: want %q, got %q", i, step, execOrder[i])
			}
		}
	})

	// Test /api/v1/users - should have authMW + loggingMW
	t.Run("users endpoint has auth and logging middleware", func(t *testing.T) {
		execOrder = []string{} // reset
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rr := httptest.NewRecorder()
		users.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		expectedOrder := []string{"auth", "logging", "list-users"}
		if len(execOrder) != len(expectedOrder) {
			t.Fatalf("execution order length: want %v, got %v", expectedOrder, execOrder)
		}
		for i, step := range expectedOrder {
			if execOrder[i] != step {
				t.Errorf("step %d: want %q, got %q", i, step, execOrder[i])
			}
		}
	})
}

// TestNestedPrefixes_MultipleConsecutive verifies multiple consecutive prefix calls
func TestNestedPrefixes_MultipleConsecutive(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Chain multiple prefixes
	nested := r.Prefix("/api").Prefix("v1").Prefix("users").Prefix("admin")

	handler := func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "nested"}
	}
	nested.Prefix("/list").GET(handler)

	// Should be accessible at /api/v1/users/admin/list
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/admin/list", nil)
	rr := httptest.NewRecorder()
	nested.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if rr.Body.String() != "nested" {
		t.Errorf("expected 'nested', got %q", rr.Body.String())
	}
}

// TestUse_MultipleCallsAccumulate verifies multiple Use() calls accumulate middleware
func TestUse_MultipleCallsAccumulate(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Multiple consecutive Use() calls
	r = r.Use(testMiddleware1)
	r = r.Use(testMiddleware2)
	r = r.Use(testMiddleware3)

	r.Prefix("/test").GET(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	want := "321"
	got := rr.Body.String()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

// TestMixedPrefixAndMiddleware verifies prefix and middleware can be chained in any order
func TestMixedPrefixAndMiddleware(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	// Mix prefix and middleware calls
	r1 := r.Prefix("/api").Use(testMiddleware1).Prefix("/v1").Use(testMiddleware2)

	r1.Prefix("/test").GET(testHandler)

	// Should be accessible at /api/v1/test with both middleware
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	rr := httptest.NewRecorder()
	r1.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	want := "21" // middleware1 + middleware2
	got := rr.Body.String()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

// TestRouteRegistration_MultipleRoutesDifferentPrefixes verifies routes under different prefixes
func TestRouteRegistration_MultipleRoutesDifferentPrefixes(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	api := r.Prefix("/api")
	admin := r.Prefix("/admin")

	api.Prefix("/users").GET(func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "api-users"}
	})

	admin.Prefix("/users").GET(func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "admin-users"}
	})

	tests := []struct {
		path string
		want string
	}{
		{"/api/users", "api-users"},
		{"/admin/users", "admin-users"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rr.Code)
			}

			if rr.Body.String() != tt.want {
				t.Errorf("want %q, got %q", tt.want, rr.Body.String())
			}
		})
	}
}

// TestUse_ReturnsNewRouter verifies that Use() returns a new router instance
func TestUse_ReturnsNewRouter(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	r2 := r.Use(testMiddleware1)
	r3 := r.Use(testMiddleware2)

	// All should be different instances
	if r == r2 {
		t.Error("r and r2 should be different instances")
	}
	if r == r3 {
		t.Error("r and r3 should be different instances")
	}
	if r2 == r3 {
		t.Error("r2 and r3 should be different instances")
	}
}

// TestPrefix_ReturnsNewRouter verifies that Prefix() returns a new router instance
func TestPrefix_ReturnsNewRouter(t *testing.T) {
	r, err := router.New()
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	r2 := r.Prefix("/api")
	r3 := r.Prefix("/admin")

	// All should be different instances
	if r == r2 {
		t.Error("r and r2 should be different instances")
	}
	if r == r3 {
		t.Error("r and r3 should be different instances")
	}
	if r2 == r3 {
		t.Error("r2 and r3 should be different instances")
	}
}
