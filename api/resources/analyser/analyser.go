package analyser

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
	h "web-analyser/util/http"
)

type Analyser struct {
	summary *Summary
	client  *http.Client
}

type inputField struct {
	t string
}

func NewAnalyser(summary *Summary) *Analyser {
	return &Analyser{
		summary: summary,
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

	a.ProcessHTML(doc)
	return nil, statusCode
}

func (a *Analyser) ProcessHTML(n *html.Node) {
	switch t := n.Type; t {
	case html.DoctypeNode:
		a.summary.SetVersion(strings.ToUpper(htmlVersion(n)))
	case html.ElementNode:
		switch tagName := n.Data; tagName {
		case "title":
			if n.FirstChild != nil {
				a.summary.SetTitle(text(n))
			}
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
			switch {
			case strings.HasPrefix(link, "//") || strings.HasPrefix(link, "http"):
				a.summary.AddExternalLink(link)
			case strings.HasPrefix(link, "/"):
				a.summary.AddInternalLink(link)
			case !strings.HasPrefix(link, "#") && !strings.HasPrefix(link, "mailto") &&
				!strings.HasPrefix(link, "tel") && !strings.HasPrefix(link, "javascript"):
				a.summary.AddInaccessibleLink(link)
			}
		case "form":
			if hasLoginForm(n) {
				a.summary.SetHasLoginForm(true)
			}

		}

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.ProcessHTML(c)
	}
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

func inputFields(n *html.Node) []*inputField {
	if n.Type == html.ElementNode && n.Data == "input" {
		field := inputField{}
		for _, attr := range n.Attr {
			if attr.Key == "type" {
				field.t = attr.Val
			}
		}
		return []*inputField{&field}
	}

	var fields []*inputField
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fields = append(fields, inputFields(c)...)
	}

	return fields
}

func hasLoginForm(n *html.Node) bool {
	actionValue := ""
	for _, attr := range n.Attr {
		if attr.Key == "action" {
			actionValue = attr.Val
			break
		}
	}
	if !strings.Contains(actionValue, "register") && !strings.Contains(actionValue, "signup") {
		if strings.Contains(actionValue, "login") || strings.Contains(actionValue, "signin") {
			return true
		} else {
			passwordInputFieldCount := 0
			usernameInputFieldCount := 0
			for _, inputField := range inputFields(n) {
				if inputField.t == "password" {
					passwordInputFieldCount++
				}
				if inputField.t == "email" || inputField.t == "text" {
					usernameInputFieldCount++
				}
			}

			if usernameInputFieldCount == 1 && passwordInputFieldCount == 1 {
				return true
			}
		}
	}
	return false
}

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
