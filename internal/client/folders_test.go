package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFolders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Folder{
			{ID: 1, UID: "f1", Title: "Folder 1"},
			{ID: 2, UID: "f2", Title: "Folder 2"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	results, err := c.ListFolders(context.Background(), PageParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(results))
	}
	if results[0].Title != "Folder 1" {
		t.Errorf("got title %q, want %q", results[0].Title, "Folder 1")
	}
}

func TestGetFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/f1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Folder{
			ID:    1,
			UID:   "f1",
			Title: "Test Folder",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetFolder(context.Background(), "f1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Title != "Test Folder" {
		t.Errorf("got title %q, want %q", result.Title, "Test Folder")
	}
}

func TestCreateFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req FolderCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "New Folder" {
			t.Errorf("got title %q, want %q", req.Title, "New Folder")
		}

		json.NewEncoder(w).Encode(Folder{
			ID:    3,
			UID:   "new-f",
			Title: "New Folder",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateFolder(context.Background(), FolderCreateRequest{Title: "New Folder"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UID != "new-f" {
		t.Errorf("got UID %q, want %q", result.UID, "new-f")
	}
}

func TestUpdateFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/f1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req FolderUpdateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "Updated Folder" {
			t.Errorf("got title %q, want %q", req.Title, "Updated Folder")
		}

		json.NewEncoder(w).Encode(Folder{
			ID:      1,
			UID:     "f1",
			Title:   "Updated Folder",
			Version: 2,
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.UpdateFolder(context.Background(), "f1", FolderUpdateRequest{Title: "Updated Folder"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version != 2 {
		t.Errorf("got version %d, want %d", result.Version, 2)
	}
}

func TestDeleteFolder(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/f1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Folder deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteFolder(context.Background(), "f1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetFolderPermissions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/f1/permissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]FolderPermission{
			{ID: 1, Role: "Viewer", Permission: 1, PermissionName: "View"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	perms, err := c.GetFolderPermissions(context.Background(), "f1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("expected 1 permission, got %d", len(perms))
	}
	if perms[0].PermissionName != "View" {
		t.Errorf("got permission name %q, want %q", perms[0].PermissionName, "View")
	}
}

func TestUpdateFolderPermissions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/folders/f1/permissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req FolderPermissionsUpdateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if len(req.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(req.Items))
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Permissions updated"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.UpdateFolderPermissions(context.Background(), "f1", FolderPermissionsUpdateRequest{
		Items: []FolderPermissionUpdate{
			{Role: "Viewer", Permission: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
