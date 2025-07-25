package types

import (
	"net/http"
	"net/url"
)

type ProblemDetails struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Detail        string         `json:"detail"`
	Instance      string         `json:"instance"`
	InvalidParams []InvalidParam `json:"invalid-params,omitempty"`
	Status        int            `json:"status"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func NewProblemDetails(r *http.Request, resourceUrl *url.URL, title, detail string, status int, params ...InvalidParam) ProblemDetails {
	return ProblemDetails{
		Type:          resourceUrl.String(),
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      r.URL.String(),
		InvalidParams: append([]InvalidParam{}, params...),
	}
}
