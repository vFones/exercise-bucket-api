package types

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProblemDetailsFromError(t *testing.T) {
	tests := []struct {
		err             error
		request         *http.Request
		typeDescription string
		name            string
		title           string
		details         string
		want            ProblemDetails
		status          int
	}{
		{
			typeDescription: "about:blank",
			name:            "Standard Error",
			err:             errors.New("an error occurred"),
			request:         httptest.NewRequest("GET", "/test-url", nil),
			title:           "Error",
			details:         "an error occurred",
			status:          500,
			want: ProblemDetails{
				Type:     "about:blank",
				Title:    "Error",
				Status:   500,
				Detail:   "an error occurred",
				Instance: "/test-url",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProblemDetails(tt.request, &url.URL{
				Scheme: "about",
				Opaque: "blank",
			}, tt.title, tt.details, tt.status)
			if got.Type != tt.want.Type || got.Title != tt.want.Title || got.Status != tt.want.Status || got.Detail != tt.want.Detail || got.Instance != tt.want.Instance {
				t.Errorf("ProblemDetails() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
