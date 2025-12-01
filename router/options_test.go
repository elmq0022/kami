package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
	"github.com/elmq0022/kami/types"
)

type testResponder struct {
	Status int
	Body   string
}

func (r *testResponder) Respond(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(r.Status)
	w.Write([]byte(r.Body))
}

func testMiddleware1(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		responder := next(r)
		return &testResponder{
			Status: responder.(*testResponder).Status,
			Body:   responder.(*testResponder).Body + "1",
		}
	}
}

func testMiddleware2(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		responder := next(r)
		return &testResponder{
			Status: responder.(*testResponder).Status,
			Body:   responder.(*testResponder).Body + "2",
		}
	}
}

func testMiddleware3(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		responder := next(r)
		return &testResponder{
			Status: responder.(*testResponder).Status,
			Body:   responder.(*testResponder).Body + "3",
		}
	}
}

func testHandler(req *http.Request) types.Responder {
	return &testResponder{Status: 200, Body: ""}
}

func TestUse(t *testing.T) {
	r, _ := router.New()
	r.Use(testMiddleware1, testMiddleware2, testMiddleware3)
	r.GET("/", testHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want %d got %d", http.StatusOK, rr.Code)
	}

	want := "321"
	got := rr.Body.String()
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestWithNotFound(t *testing.T) {
	testNotFound := func(r *http.Request) types.Responder {
		return &testResponder{
			Status: http.StatusNotFound,
			Body:   "test not found",
		}
	}

	r, _ := router.New(router.WithNotFound(testNotFound))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("want %d got %d", http.StatusNotFound, rr.Code)
	}

	if rr.Body.String() != "test not found" {
		t.Fatalf("want %s, got %s", "test not found", rr.Body.String())
	}
}

func TestLogger(t *testing.T) {
	r, _ := router.New()
	r.Use(router.Logger)
	r.GET("/test", func(req *http.Request) types.Responder {
		return &testResponder{Status: http.StatusOK, Body: "logged"}
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want %d got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "logged" {
		t.Fatalf("want %s, got %s", "logged", rr.Body.String())
	}
}

func TestRouteSpecificMiddleware(t *testing.T) {
	r, _ := router.New()
	r.Use(testMiddleware1) // Global middleware

	// Route with route-specific middleware
	r.GET("/with-mw", testHandler, testMiddleware2, testMiddleware3)

	// Route without route-specific middleware
	r.GET("/without-mw", testHandler)

	t.Run("route with middleware", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/with-mw", nil)
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("want %d got %d", http.StatusOK, rr.Code)
		}

		// Global (mw1) -> Route-specific (mw2, mw3) -> handler
		// Results in: 1 + 2 + 3 (reverse order due to wrapping)
		want := "321"
		got := rr.Body.String()
		if got != want {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	t.Run("route without middleware", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/without-mw", nil)
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("want %d got %d", http.StatusOK, rr.Code)
		}

		// Only global middleware (mw1)
		want := "1"
		got := rr.Body.String()
		if got != want {
			t.Fatalf("want %s, got %s", want, got)
		}
	})
}
