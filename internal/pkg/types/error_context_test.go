package types

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetErrorInRequestContext(t *testing.T) {
	u := url.URL{
		Scheme: "https",
		Host:   "example.com",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	customMessage := "test custom error"
	testError := errors.New("test error")
	SetErrorInRequestContext(req, testError, customMessage)
	ctx := req.Context()
	assert.ErrorIs(t, ctx.Value(ErrKey).(error), testError)
	if ctx.Value(ErrKey) != testError {
		t.Fatalf("error not inserted in context. wanted: %v, got: %v", testError, ctx.Value(ErrKey))
	}
	customErrorMessage := GetCustomErrorFromContext(ctx)
	if customErrorMessage != customMessage {
		t.Fatalf("error custom error message in context. wanted: %v, got: %v", customMessage, customErrorMessage)
	}
}
