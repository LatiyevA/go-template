package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthAndReadiness(t *testing.T) {
	for _, tc := range []struct {
		name  string
		path  string
		ready func(context.Context) error
		want  int
	}{
		{"live with unavailable database", "/healthz", func(context.Context) error { return errors.New("offline") }, 200},
		{"ready", "/readyz", func(context.Context) error { return nil }, 200},
		{"database unavailable", "/readyz", func(context.Context) error { return errors.New("offline") }, 503},
		{"missing readiness check", "/readyz", nil, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			log := slog.New(slog.NewJSONHandler(io.Discard, nil))
			h := NewHandler(nil, log, tc.ready)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d", rec.Code, tc.want)
			}
		})
	}
}

func TestReadinessPreservesRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	h := NewHandler(nil, slog.New(slog.NewJSONHandler(io.Discard, nil)), func(ctx context.Context) error {
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Error("request cancellation was lost")
		}
		if _, ok := ctx.Deadline(); !ok {
			t.Error("readiness check has no deadline")
		}
		return ctx.Err()
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil).WithContext(ctx))
	if rec.Code != 503 {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestPanicIsRecoveredAndLoggedAsJSON(t *testing.T) {
	var logs bytes.Buffer
	h := NewHandler(nil, slog.New(slog.NewJSONHandler(&logs, nil)), func(context.Context) error {
		panic("test failure")
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if rec.Code != 500 {
		t.Fatalf("status = %d", rec.Code)
	}
	decoder := json.NewDecoder(&logs)
	var recovery, access map[string]any
	if err := decoder.Decode(&recovery); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&access); err != nil {
		t.Fatal(err)
	}
	if recovery["level"] != "ERROR" || access["status"] != float64(500) {
		t.Fatalf("unexpected recovery/access logs: %v %v", recovery, access)
	}
}
