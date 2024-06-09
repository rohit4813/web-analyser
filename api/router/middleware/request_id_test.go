package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyser/api/router/middleware"
	iCtx "web-analyser/internal/utils/ctx"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
	}{
		{
			name:        "with header value",
			headerValue: "HEADER123",
		},
		{
			name:        "without header value",
			headerValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tt.headerValue != "" {
				r.Header.Set(middleware.RequestIDHeaderKey, tt.headerValue)
			}

			w := httptest.NewRecorder()
			middleware.RequestID(http.HandlerFunc(testHandlerFuncRequestID())).ServeHTTP(w, r)

			if w.Result().StatusCode != http.StatusOK {
				t.Fatal("context requestID should not be empty")
			}
		})
	}
}

func testHandlerFuncRequestID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := iCtx.RequestID(r.Context())

		if requestId == "" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
