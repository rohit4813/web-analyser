package analyser

import (
	"errors"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	url2 "net/url"
	uCtx "web-analyser/util/ctx"
	e "web-analyser/util/error"
	h "web-analyser/util/http"
)

type Handler struct {
	logger *zerolog.Logger
	tpl    *template.Template
	client *http.Client
}

func NewHandler(logger *zerolog.Logger, tpl *template.Template, client *http.Client) *Handler {
	return &Handler{
		logger: logger,
		tpl:    tpl,
		client: client,
	}
}

// Index serves the index  page.
func (a *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// get the request id from the request context
	reqID := uCtx.RequestID(r.Context())

	// render the index page
	err := a.tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		a.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
		return
	}
}

// Summary gives the summary of the url.
func (a *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	// get the url from the request
	url := r.FormValue("url")

	reqID := uCtx.RequestID(r.Context())

	// Validate the URL using the IsValidURL function
	if !h.IsValidURL(url) {
		// If the URL is invalid, log and render the error page with proper message
		a.logger.Error().Str(uCtx.KeyRequestID, reqID).Str("url", url).
			Err(errors.New("invalid URL")).Msg("")
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", e.Error{Msg: string(e.InvalidURLError)})
		if err != nil {
			a.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}

	parsedUrl, _ := url2.Parse(url)
	analyser := NewAnalyser(parsedUrl, NewSummary(), a.client)

	// analyse the page
	err, statusCode := analyser.Analyse()
	if err != nil {
		displayErr := e.Error{Msg: string(e.UnreachableURLError)}
		// if there is statusCode returned, add it to the error object so that it can be displayed in the front end
		if statusCode != nil {
			displayErr.HTTPStatusCode = *statusCode
		}
		// log and render the error page with proper message
		a.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg(displayErr.Msg)
		err := a.tpl.ExecuteTemplate(w, "error.gohtml", displayErr)
		if err != nil {
			a.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
			return
		}
		return
	}

	// render the summary page after analysing the url
	err = a.tpl.ExecuteTemplate(w, "summary.gohtml", analyser.GetSummary())
	if err != nil {
		a.logger.Error().Str(uCtx.KeyRequestID, reqID).Err(err).Msg("template error")
		return
	}
}
