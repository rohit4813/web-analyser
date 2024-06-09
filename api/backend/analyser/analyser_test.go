package analyser_test

import (
	"errors"
	"fmt"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"web-analyser/api/backend/analyser"
	"web-analyser/mocks"
)

func TestAnalyserImpl_Analyse(t *testing.T) {
	tests := []*struct {
		name                   string
		url                    string
		setupExpectations      func(client *mocks.MockClient)
		htmlData               string
		expectedSummary        *analyser.Summary
		expectedError          error
		expectedHttpStatusCode int
	}{
		{
			"Should return error if unable to reach the url",
			"https://google.com",
			func(client *mocks.MockClient) {
				client.EXPECT().Get("https://google.com").Return(nil, errors.New("error"))
			},
			"",
			nil,
			errors.New("error"),
			0,
		},
		{
			"Should return error with http status code if unable to reach the url and http status code is present",
			"https://google.com",
			func(client *mocks.MockClient) {
				resp := &http.Response{
					StatusCode: 100,
				}
				client.EXPECT().Get("https://google.com").Return(resp, errors.New("error"))
			},
			"",
			nil,
			errors.New("error"),
			100,
		},
		{
			"Should return error if status code is not 200",
			"https://google.com",
			func(client *mocks.MockClient) {
				resp := &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("")),
				}

				client.EXPECT().Get("https://google.com").Return(resp, nil)
			},
			"",
			nil,
			errors.New(fmt.Sprintf("http status code not 200, value: %v", 404)),
			404,
		},
		{
			name: "Should update the version, title and headers count in the summary",
			url:  "https://google.com",
			setupExpectations: func(client *mocks.MockClient) {
				d := "<!DOCTYPE HTML><html><title>This is the title</title><p>Halo, <b>wie geht's</b>. It means," +
					" \"Hello, <h6>how are you</h6>\".</p><h1>H1 heading.<h2>H2 heading inside H1 heading</h2></h1></html>"
				resp := &http.Response{
					Body:       io.NopCloser(strings.NewReader(d)),
					StatusCode: 200,
				}

				client.EXPECT().Get("https://google.com").Return(resp, nil)
			},
			expectedSummary: &analyser.Summary{
				Version:              "HTML 5",
				Title:                "This is the title",
				HeadersCount:         map[string]int{"h1": 1, "h2": 1, "h6": 1},
				InternalLinksMap:     map[string]struct{}{},
				ExternalLinksMap:     map[string]struct{}{},
				InaccessibleLinksMap: map[string]struct{}{},
				HasLoginForm:         false,
			},
		},
		{
			name: "Should update the internal, external and inaccessible links in the summary",
			url:  "https://google.com",
			setupExpectations: func(client *mocks.MockClient) {
				d := "<html><a href='internal_link1'></a><a href='https://google.com/internal_link2'></a>" +
					"<a href='#internal_link3'></a><a href='https://www.facebook.com/external_link1'></a>" +
					"<a href='//abc.google.com/external_link2'></a><a href='abc%$^inaccessible_link1'>" +
					"<a href='mailto:invalid_link1@abc.com'><a href='tel invalid_link2'></html>"
				resp := &http.Response{
					Body:       io.NopCloser(strings.NewReader(d)),
					StatusCode: 200,
				}

				client.EXPECT().Get("https://google.com").Return(resp, nil)
			},
			expectedSummary: &analyser.Summary{
				Version:      "",
				Title:        "",
				HeadersCount: map[string]int{},
				InternalLinksMap: map[string]struct{}{
					"internal_link1":                    struct{}{},
					"https://google.com/internal_link2": struct{}{},
					"#internal_link3":                   struct{}{},
				},
				ExternalLinksMap: map[string]struct{}{
					"https://www.facebook.com/external_link1": struct{}{},
					"//abc.google.com/external_link2":         struct{}{},
				},
				InaccessibleLinksMap: map[string]struct{}{
					"abc%$^inaccessible_link1": struct{}{},
				},
				HasLoginForm: false,
			},
		},
		{
			name: "Should not update the has login form in the summary",
			url:  "https://google.com",
			setupExpectations: func(client *mocks.MockClient) {
				d := "<!DOCTYPE HTML><form action='/signup'></html>"
				resp := &http.Response{
					Body:       io.NopCloser(strings.NewReader(d)),
					StatusCode: 200,
				}

				client.EXPECT().Get("https://google.com").Return(resp, nil)
			},
			expectedSummary: &analyser.Summary{
				Version:              "HTML 5",
				Title:                "",
				HeadersCount:         map[string]int{},
				InternalLinksMap:     map[string]struct{}{},
				ExternalLinksMap:     map[string]struct{}{},
				InaccessibleLinksMap: map[string]struct{}{},
				HasLoginForm:         false,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockClient(ctrl)
			a := analyser.NewAnalyser(mockClient)
			u := tc.url

			if tc.setupExpectations != nil {
				tc.setupExpectations(mockClient)
			}

			parsedUrl, _ := url.Parse(u)
			summary, err, statusCode := a.Analyse(parsedUrl)
			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Fatalf("Expected:%v, Got:%v", tc.expectedError, err)
			}
			if tc.expectedHttpStatusCode != 0 && statusCode != tc.expectedHttpStatusCode {
				t.Fatalf("Expected:%v, Got:%v", tc.expectedHttpStatusCode, statusCode)
			}

			if tc.expectedSummary != nil {
				tc.expectedSummary.URL = parsedUrl
				if !reflect.DeepEqual(tc.expectedSummary, summary) {
					t.Fatalf("Expected:%+v, Got:%+v", tc.expectedSummary, summary)
				}
			}

		})
	}
}
