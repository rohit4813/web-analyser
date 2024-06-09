package http_test

import (
	"testing"
	"web-analyser/internal/utils/http"
)

type testCase struct {
	url    string
	expect bool
}

var tests = []*testCase{
	{
		url:    "http://www.google.com",
		expect: true,
	},
	{
		url:    "www.google.com",
		expect: false,
	},
	{
		url:    "htp://www.google.com",
		expect: false,
	},
	{
		url:    "//www.google.com",
		expect: false,
	},
	{
		url:    "http://google.com",
		expect: true,
	},
	{
		url:    "http://google",
		expect: false,
	},
	{
		url:    "http://www.google.com/abcdef/123456",
		expect: true,
	},
	{
		url:    "http://abc.def/",
		expect: true,
	},
}

func TestIsValidURL(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			validURL := http.IsValidURL(tc.url)
			if validURL != tc.expect {
				t.Fatalf("Expected:%v, Got:%v", tc.expect, validURL)
			}
		})
	}
}
