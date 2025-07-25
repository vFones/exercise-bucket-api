package types

import (
	"context"
	"net/http"
)

func SetErrorInRequestContext(request *http.Request, err error, message string) {
	ctx := request.Context()
	errorContext := context.WithValue(ctx, ErrKey, err)
	messageContext := context.WithValue(errorContext, ErrorMessageKey, message)

	*request = *request.WithContext(messageContext)
}

func GetCustomErrorFromContext(ctx context.Context) string {
	customErr := ctx.Value(ErrorMessageKey)
	customErrMessage := ""
	if customErr != nil {
		customErrMessage = customErr.(string)
	}
	return customErrMessage
}
