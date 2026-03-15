package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressGzip(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	handler.ServeHTTP(rec, req)

	if v := rec.Header().Get("Content-Encoding"); v != "gzip" {
		t.Fatalf("expected Content-Encoding gzip, got %q", v)
	}

	gr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	body, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("failed to read gzip body: %v", err)
	}

	if string(body) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(body))
	}
}

func TestCompressNoGzipWithoutAcceptEncoding(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	if v := rec.Header().Get("Content-Encoding"); v != "" {
		t.Errorf("expected no Content-Encoding, got %q", v)
	}
	if rec.Body.String() != "hello" {
		t.Errorf("expected 'hello', got %q", rec.Body.String())
	}
}

func TestCompressSkipsCompressedContentType(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("fake png data"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	handler.ServeHTTP(rec, req)

	if rec.Body.String() != "fake png data" {
		t.Errorf("expected raw body for compressed content type, got %q", rec.Body.String())
	}
}
