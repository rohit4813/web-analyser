package analyser

import (
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

// Summary represents a summary of HTML page.
type Summary struct {
	version              string              // version represents the HTML version.
	title                string              // title represents the HTML page title.
	headerCount          map[string]int      // headerCount represents the count of each header type.
	internalLinksMap     map[string]struct{} // internalLinksMap represents internal links found in the HTML page.
	externalLinksMap     map[string]struct{} // externalLinksMap represents external links found in the HTML page.
	inaccessibleLinksMap map[string]struct{} // inaccessibleLinksMap represents inaccessible links found in the HTML page.
	hasLoginForm         bool                // hasLoginForm represents if the HTML page contains a login form.
}

// NewSummary creates a new instance of Summary
func NewSummary() *Summary {
	return &Summary{
		headerCount:          make(map[string]int),
		internalLinksMap:     map[string]struct{}{},
		externalLinksMap:     map[string]struct{}{},
		inaccessibleLinksMap: map[string]struct{}{},
	}
}

// SetVersion sets the HTML version
func (s *Summary) SetVersion(version string) {
	s.version = version
}

// SetTitle sets the HTML page title
func (s *Summary) SetTitle(title string) {
	s.title = title
}

// IncrementHeaderCount increments the count for the specified header type
func (s *Summary) IncrementHeaderCount(header string) {
	s.headerCount[header]++
}

// AddExternalLink adds a link to the externalLinks
func (s *Summary) AddExternalLink(link string) {
	if _, ok := s.externalLinksMap[link]; !ok {
		s.externalLinksMap[link] = struct{}{}
	}
}

// AddInternalLink adds a link to the internalLinks
func (s *Summary) AddInternalLink(link string) {
	if _, ok := s.internalLinksMap[link]; !ok {
		s.internalLinksMap[link] = struct{}{}
	}
}

// AddInaccessibleLink adds a link to the inaccessibleLinks
func (s *Summary) AddInaccessibleLink(link string) {
	if _, ok := s.inaccessibleLinksMap[link]; !ok {
		s.inaccessibleLinksMap[link] = struct{}{}
	}
}

// SetHasLoginForm sets the has login form flag
func (s *Summary) SetHasLoginForm(status bool) {
	s.hasLoginForm = status
}

func (s *Summary) UpdateAttributes(n *html.Node) {
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
			if hasLoginForm(n) {
				s.SetHasLoginForm(true)
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
