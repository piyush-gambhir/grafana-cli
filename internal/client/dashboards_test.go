package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchDashboards(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/api/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "dash-db" {
			t.Errorf("expected type=dash-db, got %s", r.URL.Query().Get("type"))
		}
		json.NewEncoder(w).Encode([]DashboardSearchResult{
			{ID: 1, UID: "abc", Title: "Test Dashboard"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	results, err := c.SearchDashboards(context.Background(), "", "", "", PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].UID != "abc" {
		t.Errorf("got UID %q, want %q", results[0].UID, "abc")
	}
}

func TestSearchDashboards_WithFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("query") != "prod" {
			t.Errorf("expected query=prod, got %s", q.Get("query"))
		}
		if q.Get("tag") != "monitoring" {
			t.Errorf("expected tag=monitoring, got %s", q.Get("tag"))
		}
		if q.Get("folderUIDs") != "folder-uid" {
			t.Errorf("expected folderUIDs=folder-uid, got %s", q.Get("folderUIDs"))
		}
		json.NewEncoder(w).Encode([]DashboardSearchResult{})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	_, err := c.SearchDashboards(context.Background(), "prod", "monitoring", "folder-uid", PageParams{Page: 1, PerPage: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDashboardByUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/abc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(DashboardFullResponse{
			Meta: DashboardMeta{Slug: "test-dash", Version: 3},
			Dashboard: map[string]interface{}{
				"id":    float64(42),
				"title": "Test",
			},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetDashboardByUID(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Meta.Slug != "test-dash" {
		t.Errorf("got slug %q, want %q", result.Meta.Slug, "test-dash")
	}
	if result.Meta.Version != 3 {
		t.Errorf("got version %d, want %d", result.Meta.Version, 3)
	}
}

func TestCreateDashboard(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/db" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var req DashboardCreateRequest
		json.NewDecoder(r.Body).Decode(&req)

		json.NewEncoder(w).Encode(DashboardCreateResponse{
			ID:      1,
			UID:     "new-uid",
			URL:     "/d/new-uid/test",
			Status:  "success",
			Version: 1,
			Slug:    "test",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateDashboard(context.Background(), DashboardCreateRequest{
		Dashboard: map[string]interface{}{"title": "Test"},
		FolderUID: "folder1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UID != "new-uid" {
		t.Errorf("got UID %q, want %q", result.UID, "new-uid")
	}
	if result.Status != "success" {
		t.Errorf("got status %q, want %q", result.Status, "success")
	}
}

func TestUpdateDashboard(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/db" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req DashboardCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if !req.Overwrite {
			t.Error("expected overwrite=true")
		}

		json.NewEncoder(w).Encode(DashboardCreateResponse{
			ID:      1,
			UID:     "existing-uid",
			Status:  "success",
			Version: 2,
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateDashboard(context.Background(), DashboardCreateRequest{
		Dashboard: map[string]interface{}{"title": "Updated"},
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version != 2 {
		t.Errorf("got version %d, want %d", result.Version, 2)
	}
}

func TestDeleteDashboardByUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/del-uid" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Dashboard deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteDashboard(context.Background(), "del-uid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDashboardVersions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/abc123/versions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify query params
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("start") != "5" {
			t.Errorf("expected start=5, got %s", r.URL.Query().Get("start"))
		}

		json.NewEncoder(w).Encode(DashboardVersionsResponse{
			Versions: []DashboardVersion{
				{ID: 1, Version: 1, CreatedBy: "admin", Message: "Initial"},
				{ID: 2, Version: 2, CreatedBy: "admin", Message: "Updated"},
			},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	versions, err := c.GetDashboardVersions(context.Background(), "abc123", 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
	if versions[0].Version != 1 {
		t.Errorf("got version %d, want 1", versions[0].Version)
	}
	if versions[1].Message != "Updated" {
		t.Errorf("got message %q, want %q", versions[1].Message, "Updated")
	}
}

func TestRestoreDashboard(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/abc123/restore" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req DashboardRestoreRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Version != 3 {
			t.Errorf("expected version 3, got %d", req.Version)
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Dashboard restored"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.RestoreDashboardVersion(context.Background(), "abc123", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDashboardPermissions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/dashboards/uid/abc123/permissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]DashboardPermission{
			{ID: 1, DashboardUID: "abc123", Role: "Viewer", Permission: 1, PermissionName: "View"},
			{ID: 2, DashboardUID: "abc123", Role: "Editor", Permission: 2, PermissionName: "Edit"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	perms, err := c.GetDashboardPermissions(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(perms))
	}
	if perms[0].PermissionName != "View" {
		t.Errorf("got permission name %q, want %q", perms[0].PermissionName, "View")
	}
}
