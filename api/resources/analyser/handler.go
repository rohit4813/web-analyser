package analyser

import (
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	netUrl "net/url"
	uCtx "web-analyser/util/ctx"
	e "web-analyser/util/error"
	iHttp "web-analyser/util/http"
	iTemplate "web-analyser/util/template"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	Summary(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	logger   *zerolog.Logger
	tpl      iTemplate.Template
	analyser Analyser
}

func NewHandler(logger *zerolog.Logger, tpl iTemplate.Template, analyser Analyser) *HandlerImpl {
	return &HandlerImpl{
		logger:   logger,
		tpl:      tpl,
		analyser: analyser,
	}
}

// Index serves the index  page.
func (h *HandlerImpl) Index(w http.ResponseWriter, r *http.Request) {
	// get the request id from the request context
	reqID := uCtx.RequestID(r.Context())

	// render the index page
	err := h.tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		h.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
		return
	}
}

// Summary gives the summary of the url.
func (h *HandlerImpl) Summary(w http.ResponseWriter, r *http.Request) {
	// get the url from the request
	url := r.FormValue("url")

	reqID := uCtx.RequestID(r.Context())

	// Validate the URL using the IsValidURL function
	if !iHttp.IsValidURL(url) {
		// If the URL is invalid, log and render the error page with proper message
		h.logger.Error().Str(uCtx.KeyRequestID, reqID).Str("url", url).
			Err(errors.New("invalid URL")).Msg("")
		h.renderTemplate(w, r, "error.gohtml", e.Error{Msg: string(e.InvalidURLError)})
		return
	}

	parsedUrl, _ := netUrl.Parse(url)
	err, statusCode := h.analyser.Analyse(parsedUrl)
	if err != nil {
		displayErr := e.Error{Msg: string(e.UnreachableURLError)}
		// if there is statusCode returned, add it to the error object so that it can be conveyed to the user
		if statusCode != nil {
			displayErr.HTTPStatusCode = *statusCode
		}
		// log and render the error page with proper message
		h.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg(displayErr.Msg)
		h.renderTemplate(w, r, "error.gohtml", displayErr)
		return
	}

	// render the summary after analysing the url
	h.renderTemplate(w, r, "summary.gohtml", h.analyser.GetSummary())
}

func (h *HandlerImpl) renderTemplate(w http.ResponseWriter, r *http.Request, tpl string, data any) {
	reqID := uCtx.RequestID(r.Context())
	err := h.tpl.ExecuteTemplate(w, tpl, data)
	if err != nil {
		h.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
		return
	}
}
