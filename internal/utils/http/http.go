package http

import (
	"net/http"
	"regexp"
	"web-analyser/config"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

const urlPattern = `http(s?)(:\/\/)((www\.)?)(([^.]+)\.)?([a-zA-z0-9\-_]+)(\.[a-zA-z0-9\-_]+)(\/[^\s]*)?`

// NewHttpClient returns a new http client
func NewHttpClient(c *config.ClientConf) *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
	}
}

// IsValidURL checks if the given string is a expect URL using regex
func IsValidURL(input string) bool {
	// Using a basic URL regex, examples: http://abc.def.com, http://www.abc.def.com/abc
	match, err := regexp.MatchString(urlPattern, input)
	if err != nil {
		return false
	}
	return match
}
