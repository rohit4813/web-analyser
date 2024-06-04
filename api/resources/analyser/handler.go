package analyser

import (
	"fmt"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	uCtx "web-analyser/util/ctx"
	e "web-analyser/util/error"
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
		e.ServerError(w, e.Error)
		return
	}
}

// Analyse analyses the given url.
func (a *Handler) Analyse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	reqID := uCtx.RequestID(r.Context())
	analyser := NewAnalyser()
	err := analyser.Analyse(r.FormValue("url"))
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.Error)
		return
	}
	fmt.Println(fmt.Sprintf("summary: %+v", analyser.summary))
	err = a.tpl.ExecuteTemplate(w, "summary.gohtml", analyser.summary)
	if err != nil {
		a.logger.Error().Str(logger.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.Error)
		return
	}
}
