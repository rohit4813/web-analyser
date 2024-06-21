package middleware

import (
	"github.com/google/uuid"
	"net/http"
	iCtx "web-analyser/internal/utils/ctx"
)

const RequestIDHeaderKey = "X-Request-ID"

// RequestID gets the request id passed in request header, generates a new one if not present
// and sets it in the ctx to be used by other libraries and services
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(RequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = iCtx.SetRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
