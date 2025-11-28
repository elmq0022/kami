package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/handlers"
)

func TestDefaultNotFoundHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/foo", nil)
	responder := handlers.DefaultNotFoundHandler(r)
	responder.Respond(rr, r)
	dnf := responder.(*handlers.DefaultNotFoundResponder)

	if rr.Code != dnf.Status {
		t.Fatalf("want %d, got %d", dnf.Status, rr.Code)
	}

	if rr.Body.String() != dnf.Body {
		t.Fatalf("want %s, got %s", dnf.Body, rr.Body.String())
	}
}
