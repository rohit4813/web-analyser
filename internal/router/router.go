package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"web-analyser/api/resources/analyser"
	"web-analyser/api/resources/health"
	"web-analyser/internal/router/middleware"
)

// New sets the routes using chi.Mux
func New(l *zerolog.Logger, h *analyser.Handler) *chi.Mux {
	r := chi.NewRouter()
	// using RequestID middleware to set request id in ctx
	r.Use(middleware.RequestID)

	// setting route for health api end point
	r.Get("/healthy", health.Read)

	// using NewRequestLog middleware to log important fields
	r.Method(http.MethodGet, "/", middleware.NewRequestLog(h.Index, l))
	r.Method(http.MethodPost, "/summary", middleware.NewRequestLog(h.Summary, l))

	return r
}
