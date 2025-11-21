package router_test

import (
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

func NewSpyAdapter(records *[]SpyAdapterRecord) types.Adapter {
	return func(w http.ResponseWriter, r *http.Request, h types.Handler) {
		status, body, err := h(r)

		record := SpyAdapterRecord{}
		record.Status = status
		record.Body = body
		record.Err = err
		record.Params = router.GetParams(r.Context())

		*records = append(*records, record)
	}
}

func NewTestHandler(status int, body any, err error) types.Handler {
	return func(req *http.Request) (int, any, error) {
		return status, body, err
	}
}

func TestRouter_BasicRoutes(t *testing.T) {
	want := 1

	routes := types.Routes{
		{Method: http.MethodGet, Path: "/", Handler: NewTestHandler(http.StatusOK, want, nil)},
		// {Method: http.MethodGet, Path: "/:id", Handler: paramsHandler},
	}

	records := []SpyAdapterRecord{}
	r := router.New(routes, NewSpyAdapter(&records))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if want != records[0].Body {
		t.Fatalf("want %v, got %v", want, records[0].Body)
	}

	if records[0].Err != nil {
		t.Fatalf("want nil, got %q", records[0].Err)
	}
}
