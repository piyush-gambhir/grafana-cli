package client

import (
	"context"
)

// Snapshot represents a dashboard snapshot.
type Snapshot struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	OrgID     int64  `json:"orgId"`
	UserID    int64  `json:"userId"`
	External  bool   `json:"external"`
	ExternalURL string `json:"externalUrl,omitempty"`
	Expires   string `json:"expires"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	URL       string `json:"url,omitempty"`
}

// SnapshotDetail represents a snapshot with its full dashboard data.
type SnapshotDetail struct {
	Meta      SnapshotMeta           `json:"meta"`
	Dashboard map[string]interface{} `json:"dashboard"`
}

// SnapshotMeta holds snapshot metadata.
type SnapshotMeta struct {
	IsSnapshot bool   `json:"isSnapshot"`
	Type       string `json:"type"`
	Slug       string `json:"slug"`
	Created    string `json:"created"`
	Expires    string `json:"expires"`
}

// SnapshotCreateRequest is the body for creating a snapshot.
type SnapshotCreateRequest struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	Name      string                 `json:"name,omitempty"`
	Expires   int64                  `json:"expires,omitempty"`
	External  bool                   `json:"external,omitempty"`
	Key       string                 `json:"key,omitempty"`
}

// SnapshotCreateResponse is the response from creating a snapshot.
type SnapshotCreateResponse struct {
	ID        int64  `json:"id"`
	Key       string `json:"key"`
	DeleteKey string `json:"deleteKey"`
	URL       string `json:"url"`
	DeleteURL string `json:"deleteUrl"`
}

// ListSnapshots returns all snapshots.
func (c *Client) ListSnapshots(ctx context.Context) ([]Snapshot, error) {
	var results []Snapshot
	resp, err := c.Get(ctx, "/api/dashboard/snapshots")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetSnapshot returns a snapshot by key.
func (c *Client) GetSnapshot(ctx context.Context, key string) (*SnapshotDetail, error) {
	var result SnapshotDetail
	resp, err := c.Get(ctx, "/api/snapshots/"+key)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSnapshot creates a new dashboard snapshot.
func (c *Client) CreateSnapshot(ctx context.Context, req SnapshotCreateRequest) (*SnapshotCreateResponse, error) {
	var result SnapshotCreateResponse
	resp, err := c.Post(ctx, "/api/snapshots", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteSnapshot deletes a snapshot by key.
func (c *Client) DeleteSnapshot(ctx context.Context, key string) error {
	resp, err := c.Delete(ctx, "/api/snapshots/"+key)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
