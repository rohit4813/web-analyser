package html

import (
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

const htmlVersionRegex = "(^.*) ((xhtml|html) ([a-zA-z]+ )?[0-9]*\\.?[0-9]*)(.*$)"

// Text returns the appended Text from the given html node, it recursively calls itself to build the Text
func Text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += Text(c)
	}
	return ret
}

// Version finds the html version from the given html node, returns unknown if unable to do so
func Version(n *html.Node) string {
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
			re := regexp.MustCompile(htmlVersionRegex)
			matches := re.FindStringSubmatch(strings.ToLower(publicAttrValue))
			if len(matches) >= 3 {
				version = matches[2]
			}
		}
	}
	return version
}

// HasLoginForm returns whether the give form html node has a login form or not
func HasLoginForm(n *html.Node) bool {
	// actionValue is the action attribute of the form
	actionValue := ""
	for _, attr := range n.Attr {
		if attr.Key == "action" {
			actionValue = attr.Val
			break
		}
	}
	// we check if actionValue doesn't contain register and signup, that means it can be a login form
	// then we check if the actionValue contain either login or signin
	// return true if that is the case, otherwise we fall back to checking if the form has one username and
	// one password input
	if !strings.Contains(actionValue, "register") && !strings.Contains(actionValue, "signup") {
		if strings.Contains(actionValue, "login") || strings.Contains(actionValue, "signin") {
			return true
		} else {
			passwordInputFieldCount := 0
			usernameInputFieldCount := 0
			// iterate over all the input type fields in the form
			for _, inputField := range InputFields(n) {
				if inputField == "password" {
					passwordInputFieldCount++
				}
				if inputField == "email" || inputField == "Text" {
					usernameInputFieldCount++
				}
			}

			// if we have one email or Text input type field and one password input type field
			// return true
			if usernameInputFieldCount == 1 && passwordInputFieldCount == 1 {
				return true
			}
		}
	}
	return false
}

// InputFields returns all the input tag types from a html node, it recursively calls itself to build the InputFields
func InputFields(n *html.Node) []string {
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
		fields = append(fields, InputFields(c)...)
	}

	return fields
}
