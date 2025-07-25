package middleware

import (
	"context"
	"net/http"
	"time"

	"bucket_organizer/internal/app/server/httputils"
	"bucket_organizer/pkg/logger"
	"github.com/google/uuid"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := context.WithValue(r.Context(), logger.TraceId, uuid.New().String())

		logValues := make([]interface{}, 0)
		if requestID := r.Header.Get("X-Request-Id"); requestID != "" {
			logValues = append(logValues, logger.NewLogValue("requestId", requestID))
		}
		logValues = append(logValues, logger.NewLogValue("method", r.Method))
		logValues = append(logValues, logger.NewLogValue("path", r.URL.Path))
		if r.URL.RawQuery != "" {
			logValues = append(logValues, logger.NewLogValue("query", r.URL.RawQuery))
		}
		logValues = append(logValues, logger.NewLogValue("remoteAddr", r.RemoteAddr))
		logValues = append(logValues, logger.NewLogValue("realIp", r.Header.Get("X-Real-Ip")))
		logValues = append(logValues, logger.NewLogValue("userAgent", r.UserAgent()))

		*r = *r.WithContext(ctx)

		logger.InfoNoCaller(r.Context(), "http request started", logValues...)

		next.ServeHTTP(w, r)

		logValues = append(logValues, logger.NewLogValue("secondsElapsed", time.Since(start)))
		logValues = append(logValues, logger.NewLogValue("statusCode", r.Context().Value(httputils.StatusCode)))

		msg := "http request completed"
		err := ErrorChecker(r.Context())
		if err == nil {
			logger.InfoNoCaller(ctx, msg, logValues...)
			return
		}
		logValues = append(logValues, err)
		logger.ErrorStackSkipNoCaller(ctx, msg, logValues...)
	})
}
