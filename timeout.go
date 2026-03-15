package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout returns middleware that enforces a timeout on each request.
// If the handler does not complete within the given duration, the client
// receives a 503 Service Unavailable response.
func Timeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			done := make(chan struct{})
			tw := &timeoutWriter{ResponseWriter: w}

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				// Handler completed in time.
			case <-ctx.Done():
				if !tw.wroteHeader {
					http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				}
			}
		})
	}
}

// timeoutWriter tracks whether the header has been written.
type timeoutWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.wroteHeader = true
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.wroteHeader = true
	return tw.ResponseWriter.Write(b)
}
