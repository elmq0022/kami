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
		status, body, err := h(r)
		record.Status = status
		record.Body = body
		record.Err = err
		record.Params = router.GetParams(r.Context())
	}
}

func NewTestHandler(status int, body any, err error) types.Handler {
	return func(req *http.Request) (int, any, error) {
		return status, body, err
	}
}

func TestRouter_BasicRoutes(t *testing.T) {
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
		{name: "root", method: http.MethodGet, path: "/", wantStatus: http.StatusOK, wantBody: "root", wantErr: nil, wantParams: map[string]string{}, callPath: "/"},
		{name: "root", method: http.MethodGet, path: "/book/:id", wantStatus: http.StatusOK, wantBody: "books", wantErr: nil, wantParams: map[string]string{"id": "lifeOfPi"}, callPath: "/book/lifeOfPi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := SpyAdapterRecord{}
			routes := types.Routes{{Method: tt.method, Path: tt.path, Handler: NewTestHandler(tt.wantStatus, tt.wantBody, tt.wantErr)}}
			r := router.New(routes, NewSpyAdapter(&record))

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
