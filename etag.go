package middleware

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
)

// ETag returns middleware that generates an ETag header based on a SHA-256 hash
// of the response body. If the request includes an If-None-Match header that
// matches the computed ETag, a 304 Not Modified response is returned instead.
func ETag() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply to GET and HEAD requests.
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			bw := &bufferedWriter{
				ResponseWriter: w,
				buf:            &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(bw, r)

			body := bw.buf.Bytes()
			if len(body) == 0 {
				// Write status only, no body to hash.
				w.WriteHeader(bw.statusCode)
				return
			}

			hash := sha256.Sum256(body)
			etag := fmt.Sprintf(`"%x"`, hash[:16])

			if match := r.Header.Get("If-None-Match"); match == etag {
				w.Header().Set("ETag", etag)
				w.WriteHeader(http.StatusNotModified)
				return
			}

			w.Header().Set("ETag", etag)
			w.WriteHeader(bw.statusCode)
			w.Write(body)
		})
	}
}

// bufferedWriter captures the response body in a buffer so the ETag can be
// computed before the response is sent to the client.
type bufferedWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
	wroteHeader bool
}

func (bw *bufferedWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.statusCode = code
		bw.wroteHeader = true
	}
}

func (bw *bufferedWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}
