package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
	iCtx "web-analyser/internal/utils/ctx"
)

type RequestLog struct {
	handler http.Handler
	logger  *zerolog.Logger
}

// NewRequestLog returns RequestLog middleware
func NewRequestLog(h http.HandlerFunc, l *zerolog.Logger) *RequestLog {
	return &RequestLog{
		handler: h,
		logger:  l,
	}
}

// ServeHTTP logs the important fields to help in debugging issues and then calls the handler ServeHTTP
func (h *RequestLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().
		Str("request id", iCtx.RequestID(r.Context())).
		Str("method", r.Method).
		Str("request url", r.URL.String()).
		Msg("request log")
	h.handler.ServeHTTP(w, r)
}
