package client

import (
	"context"
	"fmt"
	"net/url"
)

// Playlist represents a Grafana playlist.
type Playlist struct {
	ID       int64          `json:"id"`
	UID      string         `json:"uid"`
	Name     string         `json:"name"`
	Interval string         `json:"interval"`
	Items    []PlaylistItem `json:"items,omitempty"`
}

// PlaylistItem represents an item in a playlist.
type PlaylistItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Order int    `json:"order"`
	Title string `json:"title"`
}

// PlaylistCreateRequest is the body for creating a playlist.
type PlaylistCreateRequest struct {
	Name     string         `json:"name"`
	Interval string         `json:"interval"`
	Items    []PlaylistItem `json:"items"`
}

// PlaylistUpdateRequest is the body for updating a playlist.
type PlaylistUpdateRequest struct {
	Name     string         `json:"name"`
	Interval string         `json:"interval"`
	Items    []PlaylistItem `json:"items"`
}

// ListPlaylists returns all playlists. If limit > 0, only that many are returned.
func (c *Client) ListPlaylists(ctx context.Context, query string, limit int) ([]Playlist, error) {
	v := url.Values{}
	if query != "" {
		v.Set("query", query)
	}
	if limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", limit))
	}
	path := "/api/playlists"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}

	var results []Playlist
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetPlaylist returns a playlist by UID.
func (c *Client) GetPlaylist(ctx context.Context, uid string) (*Playlist, error) {
	var result Playlist
	resp, err := c.Get(ctx, "/api/playlists/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePlaylist creates a new playlist.
func (c *Client) CreatePlaylist(ctx context.Context, req PlaylistCreateRequest) (*Playlist, error) {
	var result Playlist
	resp, err := c.Post(ctx, "/api/playlists", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePlaylist updates an existing playlist.
func (c *Client) UpdatePlaylist(ctx context.Context, uid string, req PlaylistUpdateRequest) (*Playlist, error) {
	var result Playlist
	resp, err := c.Put(ctx, fmt.Sprintf("/api/playlists/%s", uid), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePlaylist deletes a playlist by UID.
func (c *Client) DeletePlaylist(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/playlists/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
