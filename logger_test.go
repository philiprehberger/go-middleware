package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerLogsRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	handler.ServeHTTP(rec, req)

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Error("log output should contain method GET")
	}
	if !strings.Contains(output, "/test-path") {
		t.Error("log output should contain path /test-path")
	}
	if !strings.Contains(output, "status=200") {
		t.Error("log output should contain status=200")
	}
	if !strings.Contains(output, "duration=") {
		t.Error("log output should contain duration")
	}
}

func TestLoggerCapturesStatusCode(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/missing", nil)
	handler.ServeHTTP(rec, req)

	output := buf.String()
	if !strings.Contains(output, "status=404") {
		t.Errorf("expected status=404 in log, got: %s", output)
	}
}
