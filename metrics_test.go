package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricsCapturesStatusAndDuration(t *testing.T) {
	var (
		gotMethod   string
		gotPath     string
		gotStatus   int
		gotDuration time.Duration
	)

	onRequest := func(method, path string, status int, duration time.Duration) {
		gotMethod = method
		gotPath = path
		gotStatus = status
		gotDuration = duration
	}

	handler := Metrics(onRequest)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	handler.ServeHTTP(rec, req)

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %q", gotMethod)
	}
	if gotPath != "/items" {
		t.Errorf("expected path /items, got %q", gotPath)
	}
	if gotStatus != http.StatusCreated {
		t.Errorf("expected status 201, got %d", gotStatus)
	}
	if gotDuration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestMetricsDefaultStatus(t *testing.T) {
	var gotStatus int

	onRequest := func(method, path string, status int, duration time.Duration) {
		gotStatus = status
	}

	handler := Metrics(onRequest)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	if gotStatus != http.StatusOK {
		t.Errorf("expected default status 200, got %d", gotStatus)
	}
}
