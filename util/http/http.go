package http

import (
	"net/http"
	"regexp"
	"web-analyser/config"
)

// NewHTTPClient returns a new http client
func NewHTTPClient(c *config.ClientConf) *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
	}
}

// IsValidURL checks if the given string is a valid URL using regex
func IsValidURL(input string) bool {
	// Using a basic URL regex, examples: http://abc.def.com, http://www.abc.def.com/abc
	const urlPattern = `http(s?)(:\/\/)((www\.)?)(([^.]+)\.)?([a-zA-z0-9\-_]+)(\.[a-zA-z0-9\-_]+)(\/[^\s]*)?`
	match, err := regexp.MatchString(urlPattern, input)
	if err != nil {
		return false
	}
	return match
}
