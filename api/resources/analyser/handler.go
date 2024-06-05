package analyser

import (
	"errors"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	"strconv"
	uCtx "web-analyser/util/ctx"
	h "web-analyser/util/http"
	"web-analyser/util/logger"
)

type Handler struct {
	logger *zerolog.Logger
	tpl    *template.Template
}

func NewHandler(logger *zerolog.Logger, tpl *template.Template) *Handler {
	return &Handler{
		logger: logger,
		tpl:    tpl,
	}
}

// Index serves the index  page.
func (a *Handler) Index(w http.ResponseWriter, r *http.Request) {
	reqID := uCtx.RequestID(r.Context())

	err := a.tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
		return
	}
}

// Summary gives the summary of the url.
func (a *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	reqID := uCtx.RequestID(r.Context())

	// Validate the URL using the IsValidURL function
	if !h.IsValidURL(url) {
		// If the URL is invalid, render the error template with an error message
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(errors.New("invalid URL")).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", "Invalid URL")
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}
	analyser := NewAnalyser(NewSummary(url))
	err, statusCode := analyser.Analyse(url)
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		msg := "Unable to analyse at the moment, please try again"
		if statusCode != nil && *statusCode != http.StatusOK {
			msg = "response status code not ok: " + strconv.Itoa(*statusCode)
		}
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", msg)
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}

	err = a.tpl.ExecuteTemplate(w, "summary.gohtml", analyser.GetSummary())
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
		return
	}
}
