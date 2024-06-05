package analyser

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
	h "web-analyser/util/http"
)

type Analyser struct {
	summary *Summary
	client  *http.Client
}

type inputField struct {
	t string
}

func NewAnalyser() *Analyser {
	return &Analyser{
		summary: NewSummary(),
		client:  h.NewHTTPClient(),
	}
}

func (a *Analyser) GetSummary() *Summary {
	return a.summary
}

func (a *Analyser) Analyse(u string) (error, *int) {
	resp, err := a.client.Get(u)
	var statusCode *int
	if err != nil {
		return err, statusCode
	}
	defer resp.Body.Close()

	statusCode = &resp.StatusCode
	if *statusCode != http.StatusOK {
		return errors.New("response status code not ok"), statusCode
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err, statusCode
	}

	a.summary.UpdateAttributes(doc)
	return nil, statusCode
}
