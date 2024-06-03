package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"web-analyser/internal/api/resources/analyser"
	"web-analyser/internal/api/resources/health"
	"web-analyser/internal/router/handler"
	middleware "web-analyser/internal/router/middleware"
)

func New(l *zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Get("/healthy", health.Read)

	analyse := analyser.New(l)
	r.Method(http.MethodGet, "/", handler.NewHandler(analyse.Index))
	//r.Method(http.MethodPost, "/summary", requestlog.NewHandler(bookAPI.Create, l))

	return r
}
