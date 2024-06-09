package analyser

import (
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	netUrl "net/url"
	iCtx "web-analyser/internal/utils/ctx"
	iError "web-analyser/internal/utils/error"
	iHttp "web-analyser/internal/utils/http"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	Summary(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	logger   *zerolog.Logger
	tpl      Template
	analyser Analyser
}

func NewHandler(logger *zerolog.Logger, tpl Template, analyser Analyser) *HandlerImpl {
	return &HandlerImpl{
		logger:   logger,
		tpl:      tpl,
		analyser: analyser,
	}
}

// Index serves the index  page.
func (h *HandlerImpl) Index(w http.ResponseWriter, r *http.Request) {
	// render the index page
	h.renderTemplate(w, r, "index.gohtml", nil)
}

// Summary gives the summary of the url.
func (h *HandlerImpl) Summary(w http.ResponseWriter, r *http.Request) {
	// get the url from the request
	url := r.FormValue("url")

	reqID := iCtx.RequestID(r.Context())

	// Validate the URL using the IsValidURL function
	if !iHttp.IsValidURL(url) {
		// If the URL is invalid, log and render the error page with proper message
		h.logger.Error().Str(iCtx.KeyRequestID, reqID).Str("url", url).
			Err(errors.New("invalid URL")).Msg("")
		h.renderTemplate(w, r, "error.gohtml", iError.CustomError{Message: string(iError.InvalidURLError)})
		return
	}

	parsedUrl, _ := netUrl.Parse(url)
	summary, err, statusCode := h.analyser.Analyse(parsedUrl)
	if err != nil {
		customError := iError.CustomError{Message: string(iError.UnreachableURLError)}
		// if there is http status code returned, add it to the error object so that it can be conveyed to the user
		if statusCode != 0 {
			customError.HttpStatusCode = statusCode
		}

		// log and render the error page with proper message
		h.logger.Error().Str(iCtx.KeyRequestID, reqID).Err(err).Msg(customError.Message)
		h.renderTemplate(w, r, "error.gohtml", customError)
		return
	}

	// render the summary after analysing the url
	h.renderTemplate(w, r, "summary.gohtml", summary)
}

func (h *HandlerImpl) renderTemplate(w http.ResponseWriter, r *http.Request, tpl string, data any) {
	reqID := iCtx.RequestID(r.Context())
	err := h.tpl.ExecuteTemplate(w, tpl, data)
	if err != nil {
		h.logger.Error().Str(iCtx.KeyRequestID, reqID).Err(err).Msg("template error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
