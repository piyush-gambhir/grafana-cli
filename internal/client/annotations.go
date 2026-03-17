package client

import (
	"context"
	"fmt"
	"net/url"
)

// Annotation represents a Grafana annotation.
type Annotation struct {
	ID          int64    `json:"id"`
	AlertID     int64    `json:"alertId,omitempty"`
	DashboardID int64   `json:"dashboardId,omitempty"`
	DashboardUID string  `json:"dashboardUID,omitempty"`
	PanelID     int64    `json:"panelId,omitempty"`
	UserID      int64    `json:"userId,omitempty"`
	UserLogin   string   `json:"login,omitempty"`
	UserEmail   string   `json:"email,omitempty"`
	Time        int64    `json:"time"`
	TimeEnd     int64    `json:"timeEnd,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Text        string   `json:"text"`
	Type        string   `json:"type,omitempty"`
	Created     int64    `json:"created,omitempty"`
	Updated     int64    `json:"updated,omitempty"`
}

// AnnotationCreateRequest is the body for creating an annotation.
type AnnotationCreateRequest struct {
	DashboardID  int64    `json:"dashboardId,omitempty"`
	DashboardUID string   `json:"dashboardUID,omitempty"`
	PanelID      int64    `json:"panelId,omitempty"`
	Time         int64    `json:"time,omitempty"`
	TimeEnd      int64    `json:"timeEnd,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Text         string   `json:"text"`
}

// AnnotationCreateResponse is the response from creating an annotation.
type AnnotationCreateResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// AnnotationUpdateRequest is the body for updating an annotation.
type AnnotationUpdateRequest struct {
	Time    int64    `json:"time,omitempty"`
	TimeEnd int64    `json:"timeEnd,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Text    string   `json:"text"`
}

// AnnotationTagsResult represents the result of listing annotation tags.
type AnnotationTagsResult struct {
	Result []AnnotationTag `json:"result"`
}

// AnnotationTag represents a tag with its count.
type AnnotationTag struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

// ListAnnotations returns annotations matching the given filters.
func (c *Client) ListAnnotations(ctx context.Context, dashboardID, panelID int64, from, to int64, tags []string, limit int64) ([]Annotation, error) {
	v := url.Values{}
	if dashboardID > 0 {
		v.Set("dashboardId", fmt.Sprintf("%d", dashboardID))
	}
	if panelID > 0 {
		v.Set("panelId", fmt.Sprintf("%d", panelID))
	}
	if from > 0 {
		v.Set("from", fmt.Sprintf("%d", from))
	}
	if to > 0 {
		v.Set("to", fmt.Sprintf("%d", to))
	}
	for _, tag := range tags {
		v.Add("tags", tag)
	}
	if limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", limit))
	}

	path := "/api/annotations?" + v.Encode()
	var results []Annotation
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetAnnotation returns an annotation by ID.
func (c *Client) GetAnnotation(ctx context.Context, id int64) (*Annotation, error) {
	var result Annotation
	resp, err := c.Get(ctx, fmt.Sprintf("/api/annotations/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAnnotation creates a new annotation.
func (c *Client) CreateAnnotation(ctx context.Context, req AnnotationCreateRequest) (*AnnotationCreateResponse, error) {
	var result AnnotationCreateResponse
	resp, err := c.Post(ctx, "/api/annotations", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateAnnotation updates an existing annotation.
func (c *Client) UpdateAnnotation(ctx context.Context, id int64, req AnnotationUpdateRequest) error {
	resp, err := c.Put(ctx, fmt.Sprintf("/api/annotations/%d", id), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// DeleteAnnotation deletes an annotation by ID.
func (c *Client) DeleteAnnotation(ctx context.Context, id int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/annotations/%d", id))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetAnnotationTags returns all annotation tags.
func (c *Client) GetAnnotationTags(ctx context.Context) (*AnnotationTagsResult, error) {
	var result AnnotationTagsResult
	resp, err := c.Get(ctx, "/api/annotations/tags")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
