package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/internal/config"
)

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	c := &Client{
		BaseURL:    ts.URL,
		HTTPClient: ts.Client(),
		UserAgent:  "grafana-cli/test",
	}
	return c, ts
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		rc      *config.ResolvedConfig
		wantErr bool
	}{
		{
			name: "valid config",
			rc: &config.ResolvedConfig{
				URL:   "http://localhost:3000",
				Token: "some-token",
			},
		},
		{
			name:    "missing URL",
			rc:      &config.ResolvedConfig{},
			wantErr: true,
		},
		{
			name: "trailing slash trimmed",
			rc: &config.ResolvedConfig{
				URL: "http://localhost:3000/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.rc)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.BaseURL == "" {
				t.Fatal("expected BaseURL to be set")
			}
			// Check trailing slash is trimmed
			if c.BaseURL[len(c.BaseURL)-1] == '/' {
				t.Errorf("BaseURL should not end with /, got %q", c.BaseURL)
			}
		})
	}
}

func TestBearerTokenAuth_Header(t *testing.T) {
	var gotAuth string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()
	c.Token = "my-secret-token"

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Bearer my-secret-token"
	if gotAuth != want {
		t.Errorf("got Authorization %q, want %q", gotAuth, want)
	}
}

func TestBasicAuth_Header(t *testing.T) {
	var gotUser, gotPass string
	var gotOK bool
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, gotOK = r.BasicAuth()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()
	c.Username = "admin"
	c.Password = "pass123"

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !gotOK {
		t.Fatal("expected basic auth to be present")
	}
	if gotUser != "admin" || gotPass != "pass123" {
		t.Errorf("got user=%q pass=%q, want admin/pass123", gotUser, gotPass)
	}
}

func TestOrgIDHeader(t *testing.T) {
	var gotOrgID string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotOrgID = r.Header.Get("X-Grafana-Org-Id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()
	c.OrgID = 42

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotOrgID != "42" {
		t.Errorf("got X-Grafana-Org-Id %q, want %q", gotOrgID, "42")
	}
}

func TestGet_ContentType(t *testing.T) {
	var gotContentType string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()

	_, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// GET should not set Content-Type since it has no body.
	if gotContentType != "" {
		t.Errorf("GET should not set Content-Type, got %q", gotContentType)
	}
}

func TestPost_ContentType(t *testing.T) {
	var gotContentType string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()

	_, err := c.Post(context.Background(), "/api/dashboards/db", map[string]string{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContentType != "application/json" {
		t.Errorf("POST should set Content-Type to application/json, got %q", gotContentType)
	}
}

func TestErrorParsing_400(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
	})
	defer ts.Close()

	resp, err := c.Get(context.Background(), "/api/org")
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	err = resp.JSON(nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("got status %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("got message %q, want %q", apiErr.Message, "bad request")
	}
}

func TestErrorParsing_401(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/org")
	err := resp.JSON(nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("got status %d, want 401", apiErr.StatusCode)
	}
}

func TestErrorParsing_403(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden"})
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/org")
	err := resp.JSON(nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("got status %d, want 403", apiErr.StatusCode)
	}
	if !IsForbidden(err) {
		t.Error("expected IsForbidden to return true")
	}
}

func TestErrorParsing_404(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not found"})
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/dashboards/uid/nonexistent")
	err := resp.JSON(nil)
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}

func TestErrorParsing_500(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal error"})
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/org")
	err := resp.JSON(nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("got status %d, want 500", apiErr.StatusCode)
	}
}

func TestErrorParsing_JSONBody(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Validation failed"})
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/org")
	err := resp.JSON(nil)
	apiErr := err.(*APIError)
	if apiErr.Message != "Validation failed" {
		t.Errorf("expected JSON message extraction, got %q", apiErr.Message)
	}
}

func TestErrorParsing_PlainTextBody(t *testing.T) {
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
	})
	defer ts.Close()

	resp, _ := c.Get(context.Background(), "/api/org")
	err := resp.JSON(nil)
	apiErr := err.(*APIError)
	if apiErr.Message != "Bad Gateway" {
		t.Errorf("expected plain text fallback, got %q", apiErr.Message)
	}
}

func TestRequestURL_PathJoining(t *testing.T) {
	var gotPath string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer ts.Close()

	_, err := c.Get(context.Background(), "/api/dashboards/uid/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/api/dashboards/uid/abc123" {
		t.Errorf("got path %q, want %q", gotPath, "/api/dashboards/uid/abc123")
	}
}

func TestRequestURL_QueryParams(t *testing.T) {
	var gotRawQuery string
	c, ts := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})
	defer ts.Close()

	_, err := c.Get(context.Background(), "/api/search?type=dash-db&query=test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotRawQuery != "type=dash-db&query=test" {
		t.Errorf("got query %q, want %q", gotRawQuery, "type=dash-db&query=test")
	}
}
