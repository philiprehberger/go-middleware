package middleware

import (
	"net/http"
	"time"
)

// Metrics calls onRequest after each request with method, path, status code, and duration.
func Metrics(onRequest func(method, path string, status int, duration time.Duration)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &metricsWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			onRequest(r.Method, r.URL.Path, sw.status, time.Since(start))
		})
	}
}

type metricsWriter struct {
	http.ResponseWriter
	status int
}

func (sw *metricsWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}
