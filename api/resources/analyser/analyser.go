package analyser

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	iHttp "web-analyser/util/http"
)

type Analyser interface {
	GetSummary() *Summary
	Analyse(url *url.URL) (error, *int)
}

type AnalyserImpl struct {
	summary *Summary
	client  iHttp.Client
}

func NewAnalyser(client iHttp.Client) *AnalyserImpl {
	return &AnalyserImpl{
		client: client,
	}
}

// GetSummary returns the summary on which the analyser operated
func (a *AnalyserImpl) GetSummary() *Summary {
	return a.summary
}

// Analyse analyses the url and updates the summary, returns the error and http status code in case of error
func (a *AnalyserImpl) Analyse(url *url.URL) (error, *int) {
	a.summary = NewSummary(url.String())
	// use the http client to get the url
	resp, err := a.client.Get(url.String())
	var statusCode *int
	if err != nil {
		return err, statusCode
	}
	defer resp.Body.Close()

	//if http status code is not 200, return error
	statusCode = &resp.StatusCode
	if *statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http status code not 200, got: %v", *statusCode)),
			statusCode
	}

	//parse the response body to build the html page tree
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err, statusCode
	}

	//process the html page tree
	a.processHTML(url, doc)
	return nil, statusCode
}

// processHTML iterates over all the html nodes in a dfs manner and updates the required attributes
func (a *AnalyserImpl) processHTML(url *url.URL, n *html.Node) {
	//get the type of the node
	switch t := n.Type; t {
	//if the type is html.DoctypeNode we can get the html version
	case html.DoctypeNode:
		a.summary.SetVersion(strings.ToUpper(htmlVersion(n)))
	case html.ElementNode:
		//check the tag of the html.ElementNode
		switch tagName := n.Data; tagName {
		case "title":
			a.summary.SetTitle(text(n))
		case "h1", "h2", "h3", "h4", "h5", "h6":
			a.summary.IncrementHeadersCount(tagName)
		case "a":
			link := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link = attr.Val
					break
				}
			}
			//if !strings.HasPrefix(link, "mailto") && !strings.HasPrefix(link, "tel") &&
			//	!strings.HasPrefix(link, "javascript") {
			u, err := url.Parse(link)
			if err != nil {
				a.summary.AddInaccessibleLink(link)
			} else {
				if u.Host == "" || u.Hostname() == url.Hostname() {
					a.summary.AddInternalLink(link)
				} else {
					a.summary.AddExternalLink(link)
				}
			}
			//}
		case "form":
			if hasLoginForm(n) {
				a.summary.SetHasLoginForm(true)
			}
		}

	}
	//calls itself in a dfs manner
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.processHTML(url, c)
	}
}

// text returns the joined text from a html node, it recursively calls itself to build the text
func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return ret
}

func htmlVersion(n *html.Node) string {
	version := "unknown"
	if n.Data == "html" {
		version = "html 5"
		publicAttrValue := ""
		for _, v := range n.Attr {
			if v.Key == "public" {
				publicAttrValue = v.Val
			}
		}
		if publicAttrValue != "" {
			re := regexp.MustCompile("(^.*) ((xhtml|html) [0-9]*\\.?[0-9]*)(.*$)")
			matches := re.FindStringSubmatch(strings.ToLower(publicAttrValue))
			if len(matches) >= 3 {
				version = matches[2]
			}
		}
	}
	return version
}

// hasLoginForm returns whether a html page has login form, input is a html form node
func hasLoginForm(n *html.Node) bool {
	// actionValue is the action attribute of the form
	actionValue := ""
	for _, attr := range n.Attr {
		if attr.Key == "action" {
			actionValue = attr.Val
			break
		}
	}
	// we check if actionValue doesn't contain register and signup, that means it can be a login form
	// if the previous condition satisfies then we check if the actionValue contain either login or signin
	// return true if that is the case, otherwise we check if the form has one username and
	// one password input type field
	if !strings.Contains(actionValue, "register") && !strings.Contains(actionValue, "signup") {
		if strings.Contains(actionValue, "login") || strings.Contains(actionValue, "signin") {
			return true
		} else {
			passwordInputFieldCount := 0
			usernameInputFieldCount := 0
			// iterate over all the input type fields in the form
			for _, inputField := range inputFields(n) {
				if inputField == "password" {
					passwordInputFieldCount++
				}
				if inputField == "email" || inputField == "text" {
					usernameInputFieldCount++
				}
			}

			// if we have one email or text input type field and one password input type field
			// return true
			if usernameInputFieldCount == 1 && passwordInputFieldCount == 1 {
				return true
			}
		}
	}
	return false
}

// inputFields returns all the input tag types from a html node, it recursively calls itself to build the inputFields
func inputFields(n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "input" {
		var field string
		for _, attr := range n.Attr {
			if attr.Key == "type" {
				field = attr.Val
				break
			}
		}
		return []string{field}
	}

	var fields []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fields = append(fields, inputFields(c)...)
	}

	return fields
}
