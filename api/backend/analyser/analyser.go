package analyser

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
	iHtml "web-analyser/internal/utils/html"
	iHttp "web-analyser/internal/utils/http"
)

type Analyser interface {
	Analyse(url *url.URL) (*Summary, error, int)
}

type AnalyserImpl struct {
	httpClient iHttp.Client
}

func NewAnalyser(httpClient iHttp.Client) *AnalyserImpl {
	return &AnalyserImpl{
		httpClient: httpClient,
	}
}

// Analyse analyses the url and sets the fields in the summary. In case of no error, it returns the summary, else it
// returns the error and the http status code in case of error
func (a *AnalyserImpl) Analyse(url *url.URL) (*Summary, error, int) {
	summary := NewSummary(url)
	// use the http client to get the html page
	resp, err := a.httpClient.Get(url.String())

	// set the status code if present
	var httpStatusCode int
	if resp != nil {
		httpStatusCode = resp.StatusCode
	}

	if err != nil {
		return nil, err, httpStatusCode
	}

	defer resp.Body.Close()

	// if http status code is not 200, it means the url entered is incorrect, so return error
	if httpStatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code not 200, value: %v",
			httpStatusCode)), httpStatusCode
	}

	// parse the response body to build the html page tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err, httpStatusCode
	}

	// process the html page tree
	a.processHTML(summary, doc)
	return summary, nil, httpStatusCode
}

// processHTML iterates over all the html nodes in a dfs manner and updates the required fields
func (a *AnalyserImpl) processHTML(summary *Summary, n *html.Node) {
	//get the type of the node
	switch t := n.Type; t {
	//if the type is html.DoctypeNode we can get the html version
	case html.DoctypeNode:
		summary.SetVersion(strings.ToUpper(iHtml.Version(n)))
	case html.ElementNode:
		//check the tag of the html.ElementNode
		switch tagName := n.Data; tagName {
		case "title":
			summary.SetTitle(iHtml.Text(n))
		case "h1", "h2", "h3", "h4", "h5", "h6":
			summary.IncrementHeadersCount(tagName)
		case "a":
			link := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link = attr.Val
					break
				}
			}
			if !strings.HasPrefix(link, "mailto") && !strings.HasPrefix(link, "tel") &&
				!strings.HasPrefix(link, "javascript") {
				u, err := url.Parse(link)
				if err != nil {
					summary.AddInaccessibleLink(link)
				} else {
					if u.Host == "" || u.Hostname() == summary.URL.Hostname() {
						summary.AddInternalLink(link)
					} else {
						summary.AddExternalLink(link)
					}
				}
			}
		case "form":
			if iHtml.HasLoginForm(n) {
				summary.SetHasLoginForm(true)
			}
		}

	}
	//calls itself in a dfs manner
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.processHTML(summary, c)
	}
}
