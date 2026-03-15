package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Compress returns middleware that gzip-compresses response bodies when the
// client indicates support via the Accept-Encoding header. Responses with
// content types that are already compressed (e.g., images) are skipped.
func Compress() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			cw := &compressWriter{
				ResponseWriter: w,
				gzWriter:       gz,
			}
			defer cw.Close()

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			// Remove Content-Length since it will change after compression.
			w.Header().Del("Content-Length")

			next.ServeHTTP(cw, r)
		})
	}
}

// compressWriter wraps http.ResponseWriter to write gzip-compressed data.
type compressWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	skip     bool
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	if cw.skip {
		return cw.ResponseWriter.Write(b)
	}

	// Check if content type is already compressed.
	ct := cw.ResponseWriter.Header().Get("Content-Type")
	if isCompressedContentType(ct) {
		cw.skip = true
		// Remove gzip headers since we're not compressing.
		cw.ResponseWriter.Header().Del("Content-Encoding")
		return cw.ResponseWriter.Write(b)
	}

	return cw.gzWriter.Write(b)
}

// Close flushes and closes the gzip writer.
func (cw *compressWriter) Close() error {
	if !cw.skip {
		return cw.gzWriter.Close()
	}
	return nil
}

// Unwrap returns the underlying ResponseWriter for middleware that needs
// access to the original writer (e.g., http.Flusher).
func (cw *compressWriter) Unwrap() http.ResponseWriter {
	return cw.ResponseWriter
}

var _ io.Closer = (*compressWriter)(nil)

// isCompressedContentType returns true for content types that are already
// compressed and should not be gzipped again.
func isCompressedContentType(ct string) bool {
	compressed := []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"image/webp",
		"image/avif",
		"application/zip",
		"application/gzip",
		"application/x-gzip",
		"application/x-bzip2",
		"application/x-7z-compressed",
		"application/x-rar-compressed",
	}
	ct = strings.ToLower(strings.TrimSpace(ct))
	for _, c := range compressed {
		if strings.HasPrefix(ct, c) {
			return true
		}
	}
	return false
}
