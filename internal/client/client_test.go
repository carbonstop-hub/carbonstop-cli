package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carbonstop/carbonstop-cli/internal/config"
)

func newTestConfig(serverURL string) *config.Config {
	return &config.Config{
		BaseURL: serverURL,
		APIKey:  "test-key",
		Timeout: 10,
	}
}

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("expected X-API-Key header")
		}
		if r.URL.Query().Get("k") != "v" {
			t.Errorf("expected query param k=v")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"code":200,"msg":"ok"}`))
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	status, body, err := c.Get("/test", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Errorf("expected 200, got %d", status)
	}
	if !strings.Contains(body, `"code":200`) {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected JSON content-type, got %s", ct)
		}
		w.WriteHeader(201)
		w.Write([]byte(`{"code":200}`))
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	status, _, err := c.Post("/test", []byte(`{"hello":"world"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 201 {
		t.Errorf("expected 201, got %d", status)
	}
}

func TestRetryOn5xx(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	status, _, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Errorf("expected 200 after retries, got %d", status)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryOn429(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(429)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	status, _, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Errorf("expected 200 after retry on 429, got %d", status)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(404)
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	status, _, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 404 {
		t.Errorf("expected 404, got %d", status)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retry on 4xx), got %d", attempts)
	}
}

func TestRetryExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	_, _, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error after all retries exhausted")
	}
	if !strings.Contains(err.Error(), "retries") {
		t.Errorf("expected retry exhaustion error, got: %v", err)
	}
}

func TestPathPrefix(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/test" {
			t.Errorf("expected /api/test, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := New(newTestConfig(srv.URL))
	c.Get("api/test", nil) // path without leading /
}
