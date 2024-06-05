package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	"web-analyser/api/resources/analyser"
	"web-analyser/api/resources/health"
	"web-analyser/internal/router/middleware"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("api/templates/*.gohtml"))
}

func New(l *zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	// health api end point
	r.Get("/healthy", health.Read)

	h := analyser.NewHandler(l, tpl)
	r.Method(http.MethodGet, "/", middleware.NewRequestLog(h.Index, l))
	r.Method(http.MethodPost, "/summary", middleware.NewRequestLog(h.Summary, l))

	return r
}
