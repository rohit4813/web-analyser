package http

import (
	"net/http"
	"time"
)

// NewHTTPClient returns new http client.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}
