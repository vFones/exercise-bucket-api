package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"bucket_organizer/internal/app/server/httputils"
	"bucket_organizer/internal/pkg/types"
)

func ErrorChecker(ctx context.Context) error {
	value := ctx.Value(types.ErrKey)
	if value == nil {
		return nil
	}
	return value.(error)
}

func ErrorResponder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &url.URL{
			Scheme: "about",
			Opaque: "blank",
		}
		defer func() {
			if err := recover(); err != nil {
				respErr := errors.New(err.(error).Error())
				types.SetErrorInRequestContext(r, respErr, "internal server error")
				pd := types.NewProblemDetails(r, u, http.StatusText(http.StatusInternalServerError), "internal server error", http.StatusInternalServerError)
				_ = httputils.Respond(w, r, pd.Status, pd)
			}
		}()

		next.ServeHTTP(w, r)
		err := ErrorChecker(r.Context())
		if err == nil {
			return
		}
		var pd types.ProblemDetails
		switch {
		case errors.Is(err, types.ErrNoBucketFound) || errors.Is(err, types.ErrNoObjectFound):
			pd = types.NewProblemDetails(r, u, http.StatusText(http.StatusBadRequest), err.Error(), http.StatusBadRequest)
		default:
			pd = types.NewProblemDetails(r, u, http.StatusText(http.StatusInternalServerError), err.Error(), http.StatusInternalServerError)
		}
		_ = httputils.Respond(w, r, pd.Status, pd)
	})
}
