package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"web-analyser/api/backend/analyser"
	"web-analyser/api/backend/health"
	middleware "web-analyser/api/router/middleware"
)

// New sets the routes using chi.Mux pkg
func New(l *zerolog.Logger, h analyser.Handler) *chi.Mux {
	r := chi.NewRouter()
	// using RequestID middleware to set request id in ctx
	r.Use(middleware.RequestID)

	// setting route for health api end point
	r.Get("/healthy", health.Read)

	// using NewRequestLog middleware to log important fields
	r.Method(http.MethodGet, "/", middleware.NewRequestLog(h.Index, l))
	r.Method(http.MethodPost, "/summary", middleware.NewRequestLog(h.Summary, l))

	// redirecting 404, 405 http response to index page for smooth UX,
	// not recommended for production environment
	r.NotFound(http.RedirectHandler("/", 301).ServeHTTP)
	r.MethodNotAllowed(http.RedirectHandler("/", 301).ServeHTTP)
	return r
}
