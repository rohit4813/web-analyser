package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
	uCtx "web-analyser/util/ctx"
)

type RequestLog struct {
	handler http.Handler
	logger  *zerolog.Logger
}

func NewRequestLog(h http.HandlerFunc, l *zerolog.Logger) *RequestLog {
	return &RequestLog{
		handler: h,
		logger:  l,
	}
}

func (h *RequestLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug().
		Str("request id", uCtx.RequestID(r.Context())).
		Str("method", r.Method).
		Str("request url", r.URL.String()).
		Msg("request log")
	h.handler.ServeHTTP(w, r)
}
