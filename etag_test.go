package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestETagSetsHeader(t *testing.T) {
	body := []byte("hello world")
	handler := ETag()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	etag := rec.Header().Get("ETag")
	if etag == "" {
		t.Fatal("expected ETag header to be set")
	}

	hash := sha256.Sum256(body)
	expected := fmt.Sprintf(`"%x"`, hash[:16])
	if etag != expected {
		t.Errorf("expected ETag %q, got %q", expected, etag)
	}

	if rec.Body.String() != "hello world" {
		t.Errorf("expected body 'hello world', got %q", rec.Body.String())
	}
}

func TestETag304OnMatch(t *testing.T) {
	body := []byte("hello world")
	hash := sha256.Sum256(body)
	etag := fmt.Sprintf(`"%x"`, hash[:16])

	handler := ETag()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", etag)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotModified {
		t.Errorf("expected status 304, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body on 304, got %q", rec.Body.String())
	}
}

func TestETagSkipsNonGET(t *testing.T) {
	handler := ETag()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("created"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	handler.ServeHTTP(rec, req)

	if v := rec.Header().Get("ETag"); v != "" {
		t.Errorf("expected no ETag for POST, got %q", v)
	}
}
