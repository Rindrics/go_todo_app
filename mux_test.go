package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rindrics/go_todo_app/config"
)

func TestNewMux(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("cannot get config: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		t.Fatalf("err is not nil: %v", err)
	}
	sut.ServeHTTP(w, r)
	resp := w.Result()
	defer cleanup()

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
