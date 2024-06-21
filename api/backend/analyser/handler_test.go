package analyser_test

import (
	"errors"
	"fmt"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	netUrl "net/url"
	"strings"
	"testing"
	"web-analyser/api/backend/analyser"
	iError "web-analyser/internal/utils/error"
	l "web-analyser/internal/utils/logger"
	"web-analyser/mocks"
)

func TestHandlerImpl_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedAnalyser := mocks.NewMockAnalyser(ctrl)
	mockedTemplate := mocks.NewMockTemplate(ctrl)
	logger := l.NewLogger(false)
	handler := analyser.NewHandler(logger, mockedTemplate, mockedAnalyser)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mockedTemplate.EXPECT().ExecuteTemplate(w, "index.gohtml", nil)
	handler.Index(w, r)
}

func TestHandlerImpl_Summary(t *testing.T) {
	tests := []*struct {
		name              string
		url               string
		setupExpectations func(*httptest.ResponseRecorder, *mocks.MockAnalyser, *mocks.MockTemplate)
	}{
		{
			"Should render error template for invalid url", "invalid url",
			func(w *httptest.ResponseRecorder, mockAnalyser *mocks.MockAnalyser,
				mockTemplate *mocks.MockTemplate) {
				mockTemplate.EXPECT().ExecuteTemplate(w, "error.gohtml", iError.CustomError{
					Message: string(iError.InvalidURLError)})
			},
		},
		{
			"Should render error template if analysing the url returns error", "https://google.com",
			func(w *httptest.ResponseRecorder, mockAnalyser *mocks.MockAnalyser,
				mockTemplate *mocks.MockTemplate) {
				u, _ := netUrl.Parse("https://google.com")
				mockAnalyser.EXPECT().Analyse(u).Return(nil, errors.New("error"), 0)
				mockTemplate.EXPECT().ExecuteTemplate(w, "error.gohtml", iError.CustomError{
					Message: string(iError.UnreachableURLError)})
			},
		},
		{
			"Should render error template with http status code if analysing the url returns " +
				"error with status code", "https://google.com",
			func(w *httptest.ResponseRecorder, mockAnalyser *mocks.MockAnalyser,
				mockTemplate *mocks.MockTemplate) {
				u, _ := netUrl.Parse("https://google.com")
				statusCode := 404
				mockAnalyser.EXPECT().Analyse(u).Return(nil, errors.New("error"), statusCode)
				mockTemplate.EXPECT().ExecuteTemplate(w, "error.gohtml", iError.CustomError{
					Message: string(iError.UnreachableURLError), HttpStatusCode: statusCode})
			},
		},
		{
			"Should render summary template after analysing the url", "https://google.com",
			func(w *httptest.ResponseRecorder, mockAnalyser *mocks.MockAnalyser,
				mockTemplate *mocks.MockTemplate) {
				u, _ := netUrl.Parse("https://google.com")
				summary := &analyser.Summary{
					URL:          u,
					Version:      "HTML 5",
					Title:        "Title",
					HeadersCount: nil,
					HasLoginForm: false,
				}
				mockAnalyser.EXPECT().Analyse(u).Return(summary, nil, 0)
				mockTemplate.EXPECT().ExecuteTemplate(w, "summary.gohtml", summary)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedAnalyser := mocks.NewMockAnalyser(ctrl)
			mockedTemplate := mocks.NewMockTemplate(ctrl)
			logger := l.NewLogger(false)
			handler := analyser.NewHandler(logger, mockedTemplate, mockedAnalyser)
			w := httptest.NewRecorder()
			url := tc.url
			body := strings.NewReader(fmt.Sprintf("url=%v", url))

			r := httptest.NewRequest(http.MethodPost, "/summary", body)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if tc.setupExpectations != nil {
				tc.setupExpectations(w, mockedAnalyser, mockedTemplate)
			}

			handler.Summary(w, r)
		})
	}
}
