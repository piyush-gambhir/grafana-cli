package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListDatasources(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Datasource{
			{ID: 1, UID: "ds1", Name: "Prometheus", Type: "prometheus"},
			{ID: 2, UID: "ds2", Name: "Loki", Type: "loki"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	results, err := c.ListDatasources(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 datasources, got %d", len(results))
	}
	if results[0].Name != "Prometheus" {
		t.Errorf("got name %q, want %q", results[0].Name, "Prometheus")
	}
}

func TestGetDatasourceByID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Datasource{
			ID:   1,
			UID:  "ds1",
			Name: "Prometheus",
			Type: "prometheus",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetDatasource(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UID != "ds1" {
		t.Errorf("got UID %q, want %q", result.UID, "ds1")
	}
}

func TestGetDatasourceByUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/uid/ds-uid" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Datasource{
			ID:   1,
			UID:  "ds-uid",
			Name: "InfluxDB",
			Type: "influxdb",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetDatasourceByUID(context.Background(), "ds-uid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "InfluxDB" {
		t.Errorf("got name %q, want %q", result.Name, "InfluxDB")
	}
}

func TestCreateDatasource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}

		var req DatasourceCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Name != "New DS" {
			t.Errorf("got name %q, want %q", req.Name, "New DS")
		}

		json.NewEncoder(w).Encode(DatasourceCreateResponse{
			ID:      3,
			Name:    "New DS",
			Message: "Datasource added",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateDatasource(context.Background(), DatasourceCreateRequest{
		Name:   "New DS",
		Type:   "prometheus",
		Access: "proxy",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 3 {
		t.Errorf("got ID %d, want %d", result.ID, 3)
	}
	if result.Message != "Datasource added" {
		t.Errorf("got message %q, want %q", result.Message, "Datasource added")
	}
}

func TestUpdateDatasource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(DatasourceCreateResponse{
			ID:      1,
			Name:    "Updated DS",
			Message: "Datasource updated",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.UpdateDatasource(context.Background(), 1, DatasourceCreateRequest{
		Name:   "Updated DS",
		Type:   "prometheus",
		Access: "proxy",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated DS" {
		t.Errorf("got name %q, want %q", result.Name, "Updated DS")
	}
}

func TestDeleteDatasourceByID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Data source deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteDatasource(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteDatasourceByUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/datasources/uid/ds-uid" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Data source deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteDatasourceByUID(context.Background(), "ds-uid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDatasourceByName(t *testing.T) {
	// The Grafana API supports GET /api/datasources/name/:name.
	// Our client doesn't have a dedicated method, so we test GetDatasourceByUID
	// with a 404 scenario to verify error handling.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Data source not found"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	_, err := c.GetDatasourceByUID(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent datasource")
	}
	if !IsNotFound(err) {
		t.Errorf("expected 404 error, got: %v", err)
	}
}
