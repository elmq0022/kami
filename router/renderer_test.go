package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/router"
)

func TestRender(t *testing.T) {
	wantStatus := http.StatusOK
	data := map[string]string{"hello": "world"}
	wantBody, err := json.Marshal(data)
	if err != nil {
		t.Fatal("could not marshal test data")
	}

	jr := router.JsonResponse{Status: wantStatus, Data: data}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.DefaultRenderer(w, r, jr)

	if w.Code != wantStatus {
		t.Fatalf("want %d, got %d", wantStatus, w.Code)
	}

	if w.Body.String() != string(wantBody) {
		t.Fatalf("want %s, got %s", string(wantBody), w.Body.String())
	}

}
