package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVerboseTransport_RedactsAuthHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "session=secret123")
		w.Header().Set("X-Request-Id", "req-abc-123")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL:    ts.URL,
		HTTPClient: ts.Client(),
		Token:      "super-secret-token",
		UserAgent:  "grafana-cli/test",
	}

	var logBuf bytes.Buffer
	c.EnableVerboseLogging(&logBuf)

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := logBuf.String()

	// Authorization header should be redacted.
	if strings.Contains(output, "super-secret-token") {
		t.Errorf("verbose output should NOT contain the actual token, got:\n%s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("verbose output should contain [REDACTED], got:\n%s", output)
	}

	// Non-auth headers should be visible.
	if !strings.Contains(output, "X-Request-Id") {
		t.Errorf("verbose output should show X-Request-Id header, got:\n%s", output)
	}

	// Set-Cookie in response should be redacted.
	if strings.Contains(output, "secret123") {
		t.Errorf("verbose output should NOT contain Set-Cookie value, got:\n%s", output)
	}
}

func TestVerboseTransport_NonAuthHeadersVisible(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Ratelimit-Remaining", "99")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL:    ts.URL,
		HTTPClient: ts.Client(),
		UserAgent:  "grafana-cli/test",
		OrgID:      42,
	}

	var logBuf bytes.Buffer
	c.EnableVerboseLogging(&logBuf)

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := logBuf.String()

	// These debugging headers should be visible.
	if !strings.Contains(output, "Content-Type") {
		t.Errorf("verbose output should show Content-Type header, got:\n%s", output)
	}
	if !strings.Contains(output, "X-Ratelimit-Remaining") {
		t.Errorf("verbose output should show X-Ratelimit-Remaining header, got:\n%s", output)
	}
}

func TestRedactAuthHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer secret-token")
	h.Set("Cookie", "session=abc")
	h.Set("Set-Cookie", "session=def")
	h.Set("X-Api-Key", "key-123")
	h.Set("Content-Type", "application/json")
	h.Set("X-Request-Id", "req-456")
	h.Set("X-Grafana-Org-Id", "42")

	redacted := redactAuthHeaders(h)

	// Auth headers should be redacted.
	for _, key := range []string{"Authorization", "Cookie", "Set-Cookie", "X-Api-Key"} {
		vals := redacted.Values(key)
		if len(vals) != 1 || vals[0] != "[REDACTED]" {
			t.Errorf("expected %s to be [REDACTED], got %v", key, vals)
		}
	}

	// Non-auth headers should be preserved.
	if redacted.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type should be preserved, got %q", redacted.Get("Content-Type"))
	}
	if redacted.Get("X-Request-Id") != "req-456" {
		t.Errorf("X-Request-Id should be preserved, got %q", redacted.Get("X-Request-Id"))
	}
	if redacted.Get("X-Grafana-Org-Id") != "42" {
		t.Errorf("X-Grafana-Org-Id should be preserved, got %q", redacted.Get("X-Grafana-Org-Id"))
	}

	// Original headers should not be modified.
	if h.Get("Authorization") != "Bearer secret-token" {
		t.Error("original Authorization header was modified")
	}
}

func TestVerboseTransport_CookieHeaderRedacted(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL:    ts.URL,
		HTTPClient: ts.Client(),
		UserAgent:  "grafana-cli/test",
	}

	var logBuf bytes.Buffer
	c.EnableVerboseLogging(&logBuf)

	// Manually create request with Cookie header.
	req, _ := http.NewRequest("GET", ts.URL+"/api/test", nil)
	req.Header.Set("Cookie", "grafana_session=mysecretcookie")
	c.HTTPClient.Do(req)

	output := logBuf.String()
	if strings.Contains(output, "mysecretcookie") {
		t.Errorf("verbose output should NOT contain cookie value, got:\n%s", output)
	}
}
