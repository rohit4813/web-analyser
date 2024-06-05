package http

import (
	"net/http"
	"regexp"
	"time"
)

// NewHTTPClient returns new http client.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

// IsValidURL checks if the given string is a valid URL using regex.
func IsValidURL(input string) bool {
	const urlPattern = `http(s?)(:\/\/)((www\.)?)(([^.]+)\.)?([a-zA-z0-9\-_]+)(\.[a-zA-z0-9\-_]+)(\/[^\s]*)?`
	match, err := regexp.MatchString(urlPattern, input)
	if err != nil {
		return false
	}
	return match
}
