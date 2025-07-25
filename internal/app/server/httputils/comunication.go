package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type AppContext string

var StatusCode = AppContext("statusCode")

func WriteHeaderAndContext(w http.ResponseWriter, statusCode int, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	ctx := context.WithValue(r.Context(), StatusCode, statusCode)
	*r = *r.WithContext(ctx)
	w.WriteHeader(statusCode)
}

func Respond[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	WriteHeaderAndContext(w, status, r)
	if !reflect.ValueOf(v).IsZero() {
		if err := json.NewEncoder(w).Encode(&v); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	}
	return nil
}
