package internalhttp

import (
	"net/http"
	"time"

	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default
		}
		next.ServeHTTP(recorder, r)
		logger.LogHTTP(r, startTime, recorder.statusCode)
	})
}
