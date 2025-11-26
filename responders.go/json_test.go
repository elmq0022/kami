package responders

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponder(t *testing.T) {
	tests := []struct {
		name           string
		responder      *JSONResponder
		expectedStatus int
		expectedBody   string
		expectedCT     string
	}{
		{
			name: "simple struct with default status",
			responder: &JSONResponder{
				Body: map[string]string{"message": "hello"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"hello"}`,
			expectedCT:     "application/json",
		},
		{
			name: "simple struct with custom status",
			responder: &JSONResponder{
				Body:   map[string]int{"count": 42},
				Status: http.StatusCreated,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"count":42}`,
			expectedCT:     "application/json",
		},
		{
			name: "array body",
			responder: &JSONResponder{
				Body:   []string{"foo", "bar", "baz"},
				Status: http.StatusOK,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `["foo","bar","baz"]`,
			expectedCT:     "application/json",
		},
		{
			name: "null body",
			responder: &JSONResponder{
				Body: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `null`,
			expectedCT:     "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.responder.Respond(w)

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
	responder := &JSONResponder{
		Body: make(chan int), // channels are not JSON marshalable
	}

	responder.Respond(w)

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
		responder      *JSONErrorResponder
		expectedStatus int
		expectedBody   string
		expectedCT     string
	}{
		{
			name: "not found error",
			responder: &JSONErrorResponder{
				Status: http.StatusNotFound,
				Msg:    "resource not found",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"msg":"resource not found"}`,
			expectedCT:     "application/problem+json",
		},
		{
			name: "bad request error",
			responder: &JSONErrorResponder{
				Status: http.StatusBadRequest,
				Msg:    "invalid input",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"msg":"invalid input"}`,
			expectedCT:     "application/problem+json",
		},
		{
			name: "internal server error",
			responder: &JSONErrorResponder{
				Status: http.StatusInternalServerError,
				Msg:    "something went wrong",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"msg":"something went wrong"}`,
			expectedCT:     "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.responder.Respond(w)

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
