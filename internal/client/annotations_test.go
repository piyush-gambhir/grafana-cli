package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAnnotations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Annotation{
			{ID: 1, Text: "Deployed v1.0", Tags: []string{"deploy"}},
			{ID: 2, Text: "Deployed v1.1", Tags: []string{"deploy"}},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	results, err := c.ListAnnotations(context.Background(), 0, 0, 0, 0, nil, 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(results))
	}
	if results[0].Text != "Deployed v1.0" {
		t.Errorf("got text %q, want %q", results[0].Text, "Deployed v1.0")
	}
}

func TestListAnnotations_WithFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("dashboardId") != "42" {
			t.Errorf("expected dashboardId=42, got %s", q.Get("dashboardId"))
		}
		if q.Get("panelId") != "5" {
			t.Errorf("expected panelId=5, got %s", q.Get("panelId"))
		}
		if q.Get("from") != "1000" {
			t.Errorf("expected from=1000, got %s", q.Get("from"))
		}
		if q.Get("to") != "2000" {
			t.Errorf("expected to=2000, got %s", q.Get("to"))
		}
		tags := q["tags"]
		if len(tags) != 2 || tags[0] != "deploy" || tags[1] != "production" {
			t.Errorf("expected tags=[deploy,production], got %v", tags)
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("type") != "annotation" {
			t.Errorf("expected type=annotation, got %s", q.Get("type"))
		}
		json.NewEncoder(w).Encode([]Annotation{})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	_, err := c.ListAnnotations(context.Background(), 42, 5, 1000, 2000, []string{"deploy", "production"}, 50, "annotation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAnnotation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Annotation{
			ID:   42,
			Text: "Incident started",
			Tags: []string{"incident"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetAnnotation(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("got ID %d, want %d", result.ID, 42)
	}
	if result.Text != "Incident started" {
		t.Errorf("got text %q, want %q", result.Text, "Incident started")
	}
}

func TestCreateAnnotation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req AnnotationCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Text != "Deploy v2.0" {
			t.Errorf("got text %q, want %q", req.Text, "Deploy v2.0")
		}

		json.NewEncoder(w).Encode(AnnotationCreateResponse{
			ID:      100,
			Message: "Annotation added",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateAnnotation(context.Background(), AnnotationCreateRequest{
		Text: "Deploy v2.0",
		Tags: []string{"deploy"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 100 {
		t.Errorf("got ID %d, want %d", result.ID, 100)
	}
	if result.Message != "Annotation added" {
		t.Errorf("got message %q, want %q", result.Message, "Annotation added")
	}
}

func TestDeleteAnnotation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Annotation deleted"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteAnnotation(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAnnotationTags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/annotations/tags" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify query params for tag filter and limit
		if r.URL.Query().Get("tag") != "deploy" {
			t.Errorf("expected tag=deploy, got %s", r.URL.Query().Get("tag"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		json.NewEncoder(w).Encode(AnnotationTagsResult{
			Result: []AnnotationTag{
				{Tag: "deploy", Count: 15},
				{Tag: "deploy:production", Count: 8},
			},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.GetAnnotationTags(context.Background(), "deploy", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Result) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result.Result))
	}
	if result.Result[0].Tag != "deploy" {
		t.Errorf("got tag %q, want %q", result.Result[0].Tag, "deploy")
	}
	if result.Result[0].Count != 15 {
		t.Errorf("got count %d, want %d", result.Result[0].Count, 15)
	}
}
