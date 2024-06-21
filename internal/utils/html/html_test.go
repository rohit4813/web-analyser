package html_test

import (
	"go.uber.org/mock/gomock"
	"golang.org/x/net/html"
	"strings"
	"testing"
	iHtml "web-analyser/internal/utils/html"
)

func TestText(t *testing.T) {
	tests := []*struct {
		name         string
		htmlData     string
		expectedText string
	}{
		{
			name:         "Should return the text of the html node",
			htmlData:     "<html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedText: "Halo, wie geht's. It means, \"Hello, how are you\".",
		},
		{
			name:         "Should return empty if no text is present",
			htmlData:     "<html><p></p></html>",
			expectedText: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doc, _ := html.Parse(strings.NewReader(tc.htmlData))
			data := iHtml.Text(doc.FirstChild)
			if data != tc.expectedText {
				t.Fatalf("Expected:%v, Got:%v", tc.expectedText, data)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []*struct {
		name            string
		htmlData        string
		expectedVersion string
	}{
		{
			name:            "Should return unknown if doctype is improper",
			htmlData:        "<!DOCTYPE test><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "unknown",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE HTML><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "html 5",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\"\n   \"http://www.w3.org/TR/html4/loose.dtd\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "html 4.01",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\"\n   \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "xhtml 1.0",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \n   \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "xhtml 1.1",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\"\n    \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "xhtml basic 1.1",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//IETF//DTD HTML 2.0//EN\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "html 2.0",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "html 3.2",
		},
		{
			name:            "Should return the version of the html page",
			htmlData:        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.0//EN\"\n    \"http://www.w3.org/TR/xhtml-basic/xhtml-basic10.dtd\"><html><p>Halo, <b>wie geht's</b>. It means, \"Hello, how are you\".</p></html>",
			expectedVersion: "xhtml basic 1.0",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doc, _ := html.Parse(strings.NewReader(tc.htmlData))
			data := iHtml.Version(doc.FirstChild)
			if data != tc.expectedVersion {
				t.Fatalf("Expected:%v, Got:%v", tc.expectedVersion, data)
			}
		})
	}
}

func TestHasLoginForm(t *testing.T) {
	tests := []*struct {
		name           string
		htmlData       string
		expectedResult bool
	}{
		{
			name:           "Should return true if form action has login",
			htmlData:       "<html><form action='/login'></form></html>",
			expectedResult: true,
		},
		{
			name:           "Should return true if form action has signin",
			htmlData:       "<html><form action='/signin'></form></html>",
			expectedResult: true,
		},
		{
			name:           "Should return false if form action has register",
			htmlData:       "<html><form action='/register/test'></form></html>",
			expectedResult: false,
		},
		{
			name:           "Should return false if form action signup",
			htmlData:       "<html><form action='/signup/123'></form></html>",
			expectedResult: false,
		},
		{
			name:           "Should return true if form has 1 email field and 1 password field",
			htmlData:       "<html><form><input type='email'></input><input type='password'></input></form></html>",
			expectedResult: true,
		},
		{
			name:           "Should return false if form has 2 password fields",
			htmlData:       "<html><form><input type='password'></input><input type='password'></input></form></html>",
			expectedResult: false,
		},
		{
			name:           "Should return false if form has 2 text fields",
			htmlData:       "<html><form><input type='text'></input><input type='text'></input></form></html>",
			expectedResult: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			doc, _ := html.Parse(strings.NewReader(tc.htmlData))
			formNode := findFormNode(doc.FirstChild)
			result := iHtml.HasLoginForm(formNode)
			if result != tc.expectedResult {
				t.Fatalf("Expected:%v, Got:%v", tc.expectedResult, result)
			}
		})
	}
}

// findFormNode finds the first form node
func findFormNode(n *html.Node) *html.Node {
	if n.Type != html.ElementNode {
		return nil
	}

	var formNode *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "form" {
			formNode = c
			return formNode
		}
		formNode = findFormNode(c)
		if formNode != nil {
			return formNode
		}

	}
	return formNode
}
