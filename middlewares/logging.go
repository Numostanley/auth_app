package middlewares

import (
	"log"
	"net/http"
	"time"
)

type StatusCapturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *StatusCapturingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *StatusCapturingResponseWriter) GetStatus() int {
	return w.statusCode
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wrappedResponseWriter := &StatusCapturingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(wrappedResponseWriter, r)

		statusCode := wrappedResponseWriter.GetStatus()

		log.Printf("[%s] %s - %s - %s - %d - %v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), statusCode, time.Since(startTime))
	})
}
