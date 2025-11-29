package responders_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elmq0022/kami/responders"
	"github.com/elmq0022/kami/types"
)

func TestJSONResponder(t *testing.T) {
	tests := []struct {
		name           string
		responder      types.Responder
		expectedStatus int
		expectedBody   string
		expectedCT     string
	}{
		{
			name:           "simple struct with default status",
			responder:      responders.JSONResponse(map[string]string{"message": "hello"}, 200),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"hello"}`,
			expectedCT:     "application/json",
		},
		{
			name:           "simple struct with custom status",
			responder:      responders.JSONResponse(map[string]int{"count": 42}, http.StatusCreated),
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"count":42}`,
			expectedCT:     "application/json",
		},
		{
			name:           "array body",
			responder:      responders.JSONResponse([]string{"foo", "bar", "baz"}, http.StatusOK),
			expectedStatus: http.StatusOK,
			expectedBody:   `["foo","bar","baz"]`,
			expectedCT:     "application/json",
		},
		{
			name:           "null body",
			responder:      responders.JSONResponse(nil, 0),
			expectedStatus: http.StatusOK,
			expectedBody:   `null`,
			expectedCT:     "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.responder.Respond(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if got := w.Header().Get("Content-Type"); got != tt.expectedCT {
				t.Errorf("expected Content-Type %q, got %q", tt.expectedCT, got)
			}

			if got := w.Body.String(); got != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, got)
			}
		})
	}
}

func TestJSONResponder_UnmarshalableData(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	responder := responders.JSONResponse(make(chan int), 0) // channels are not JSON marshalable

	responder.Respond(w, r)

	// Should silently fail - no writes to response
	if w.Code != http.StatusOK {
		t.Errorf("expected default status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("expected empty body on marshal error, got %q", w.Body.String())
	}
}

func TestJSONErrorResponder(t *testing.T) {
	tests := []struct {
		name           string
		responder      types.Responder
		expectedStatus int
		expectedBody   string
		expectedCT     string
	}{
		{
			name:           "not found error",
			responder:      responders.JSONErrorResponse("resource not found", http.StatusNotFound),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"msg":"resource not found"}`,
			expectedCT:     "application/problem+json",
		},
		{
			name:           "bad request error",
			responder:      responders.JSONErrorResponse("invalid input", http.StatusBadRequest),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"msg":"invalid input"}`,
			expectedCT:     "application/problem+json",
		},
		{
			name:           "internal server error",
			responder:      responders.JSONErrorResponse("something went wrong", http.StatusInternalServerError),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"msg":"something went wrong"}`,
			expectedCT:     "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.responder.Respond(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if got := w.Header().Get("Content-Type"); got != tt.expectedCT {
				t.Errorf("expected Content-Type %q, got %q", tt.expectedCT, got)
			}

			if got := w.Body.String(); got != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, got)
			}
		})
	}
}
