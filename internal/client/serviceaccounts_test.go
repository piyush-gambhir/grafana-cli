package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchServiceAccounts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "ci" {
			t.Errorf("expected query=ci, got %s", r.URL.Query().Get("query"))
		}

		json.NewEncoder(w).Encode(ServiceAccountSearchResult{
			TotalCount: 1,
			Page:       1,
			PerPage:    100,
			ServiceAccounts: []ServiceAccount{
				{ID: 1, Name: "ci-bot", Role: "Editor"},
			},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.ListServiceAccounts(context.Background(), "ci", PageParams{Page: 1, PerPage: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalCount != 1 {
		t.Errorf("got total count %d, want 1", result.TotalCount)
	}
	if len(result.ServiceAccounts) != 1 {
		t.Fatalf("expected 1 service account, got %d", len(result.ServiceAccounts))
	}
	if result.ServiceAccounts[0].Name != "ci-bot" {
		t.Errorf("got name %q, want %q", result.ServiceAccounts[0].Name, "ci-bot")
	}
}

func TestCreateServiceAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ServiceAccountCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Name != "deploy-bot" {
			t.Errorf("got name %q, want %q", req.Name, "deploy-bot")
		}

		json.NewEncoder(w).Encode(ServiceAccount{
			ID:   10,
			Name: "deploy-bot",
			Role: "Editor",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateServiceAccount(context.Background(), ServiceAccountCreateRequest{
		Name: "deploy-bot",
		Role: "Editor",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 10 {
		t.Errorf("got ID %d, want %d", result.ID, 10)
	}
}

func TestGetServiceAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ServiceAccount{
			ID:     10,
			Name:   "ci-bot",
			Role:   "Editor",
			Tokens: 2,
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetServiceAccount(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Tokens != 2 {
		t.Errorf("got tokens %d, want %d", result.Tokens, 2)
	}
}

func TestDeleteServiceAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Service account deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteServiceAccount(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateServiceAccountToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts/10/tokens" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ServiceAccountTokenCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Name != "ci-token" {
			t.Errorf("got name %q, want %q", req.Name, "ci-token")
		}

		json.NewEncoder(w).Encode(ServiceAccountToken{
			ID:   1,
			Name: "ci-token",
			Key:  "glsa_xxxxxxxxxxxx",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateServiceAccountToken(context.Background(), 10, ServiceAccountTokenCreateRequest{
		Name: "ci-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key == "" {
		t.Error("expected key to be set")
	}
	if result.Name != "ci-token" {
		t.Errorf("got name %q, want %q", result.Name, "ci-token")
	}
}

func TestListServiceAccountTokens(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/serviceaccounts/10/tokens" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]ServiceAccountToken{
			{ID: 1, Name: "token-1"},
			{ID: 2, Name: "token-2", HasExpired: true},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	tokens, err := c.ListServiceAccountTokens(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if !tokens[1].HasExpired {
		t.Error("expected token-2 to be expired")
	}
}
