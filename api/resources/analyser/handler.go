package analyser

import (
	"errors"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	"regexp"
	uCtx "web-analyser/util/ctx"
	"web-analyser/util/logger"
)

type Handler struct {
	logger *zerolog.Logger
	tpl    *template.Template
}

func New(logger *zerolog.Logger, tpl *template.Template) *Handler {
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
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", "Error serving the index page")
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}
}

// Summary gives the summary of the url.
func (a *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	url := r.FormValue("url")
	reqID := uCtx.RequestID(r.Context())

	// Validate the URL using the IsValidURL function
	if !IsValidURL(url) {
		// If the URL is invalid, render the error template with an error message
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(errors.New("invalid URL")).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", "Invalid URL")
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}
	analyser := NewAnalyser()
	err := analyser.Analyse(url)
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", "Unable to analyse at the moment, please try again")
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}

	err = a.tpl.ExecuteTemplate(w, "summary.gohtml", analyser.summary)
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", "Unable to analyse at the moment, please try again")
		if err != nil {
			a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}
}

// IsValidURL checks if the given string is a valid URL using regex.
func IsValidURL(input string) bool {
	const urlPattern = `http(s?)(:\/\/)((www\.)?)(([^.]+)\.)?([a-zA-z0-9\-_]+)(\.[a-zA-z0-9\-_]+)(\/[^\s]*)?`
	match, err := regexp.MatchString(urlPattern, input)
	if err != nil {
		return false
	}
	return match
}
