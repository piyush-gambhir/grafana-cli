package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListLibraryElements_QueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("perPage") != "20" || q.Get("page") != "2" || q.Has("perpage") {
			t.Errorf("unexpected pagination query: %s", r.URL.RawQuery)
		}
		if q.Get("folderFilterUIDs") != "folder-uid" || q.Has("folderFilter") {
			t.Errorf("unexpected folder filter query: %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(LibraryElementSearchResult{})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	if _, err := c.ListLibraryElements(context.Background(), "", 0, "folder-uid", PageParams{Page: 2, PerPage: 20}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
