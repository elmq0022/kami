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

func (r *testResponder) Respond(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	w.Write([]byte(r.Body))
}

// func newTestMiddleware(val int) types.Middleware {
// 	return func(h types.Handler) types.Handler {
// 		return func(req *http.Request) types.Responder {

// 		}
// 	}
// }

func testMiddleWare(next types.Handler) types.Handler {
	return func(r *http.Request) types.Responder {
		response := next(r)
		return response
	}
}

func testHandler(req *http.Request) types.Responder {
	return &testResponder{Status: 200, Body: ""}
}

func TestWithMiddleware(t *testing.T) {
	r, _ := router.New(
	// router.WithMiddleware(newTestMiddleware(1)),
	// router.WithMiddleware(newTestMiddleware(2)),
	// router.WithMiddleware(newTestMiddleware(3)),
	)
	r.GET("/", testHandler)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want %d got %d", http.StatusNotFound, rr.Code)
	}

	want := ""
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
