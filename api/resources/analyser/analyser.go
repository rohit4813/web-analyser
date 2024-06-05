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

func (s *Summary) UpdateAttributes(n *html.Node) {
	//hasLoginForm := false
	switch t := n.Type; t {
	case html.DoctypeNode:
		s.SetVersion(htmlVersion(n))
	case html.ElementNode:
		switch tagName := n.Data; tagName {
		case "title":
			if n.FirstChild != nil {
				s.SetTitle(n.FirstChild.Data)
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			s.IncrementHeaderCount(tagName)
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
				s.AddExternalLink(link)
			case strings.HasPrefix(link, "/"):
				s.AddInternalLink(link)
			case !strings.HasPrefix(link, "#"): // Considering mailto and tel as inaccessible links
				s.AddInaccessibleLink(link)
			}
		case "form":
			actionValue := ""
			for _, attr := range n.Attr {
				if attr.Key == "action" {
					actionValue = attr.Val
					break
				}
			}
			if !strings.Contains(actionValue, "register") && !strings.Contains(actionValue, "signup") {
				if strings.Contains(actionValue, "login") || strings.Contains(actionValue, "signin") {
					s.SetHasLoginForm(true)
				} else {
					passwordInputFieldCount := 0
					usernameInputFieldCount := 0
					for _, inputField := range formInputFields(n) {
						if inputField.t == "password" {
							passwordInputFieldCount++
						}
						if inputField.t == "email" || inputField.t == "text" {
							usernameInputFieldCount++
						}
					}

					if usernameInputFieldCount == 1 && passwordInputFieldCount == 1 {
						s.SetHasLoginForm(true)
					}
				}
			}

		}

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.UpdateAttributes(c)
	}
}

func htmlVersion(n *html.Node) string {
	version := "unknown"
	if n.Data == "html" {
		version = "html 5"
		publicKeyValue := ""
		for _, v := range n.Attr {
			if v.Key == "public" {
				publicKeyValue = v.Val
			}
		}
		if publicKeyValue != "" {
			re := regexp.MustCompile("(^.*) ((xhtml|html) [0-9]*\\.?[0-9]*)(.*$)")
			matches := re.FindStringSubmatch(strings.ToLower(publicKeyValue))
			if len(matches) >= 3 {
				version = matches[2]
			}
		}
	}
	return version
}

func formInputFields(n *html.Node) []*inputField {
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
		fields = append(fields, formInputFields(c)...)
	}

	return fields
}

type inputField struct {
	t string
}
