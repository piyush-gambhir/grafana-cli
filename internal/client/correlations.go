package client

import (
	"context"
	"fmt"
)

// Correlation represents a Grafana correlation.
type Correlation struct {
	UID           string            `json:"uid"`
	SourceUID     string            `json:"sourceUID"`
	TargetUID     string            `json:"targetUID"`
	Label         string            `json:"label"`
	Description   string            `json:"description"`
	Config        CorrelationConfig `json:"config"`
	Provisioned   bool              `json:"provisioned,omitempty"`
}

// CorrelationConfig holds the configuration for a correlation.
type CorrelationConfig struct {
	Type        string            `json:"type"`
	Field       string            `json:"field"`
	Target      map[string]interface{} `json:"target,omitempty"`
	Transformations []CorrelationTransformation `json:"transformations,omitempty"`
}

// CorrelationTransformation defines a transformation for a correlation.
type CorrelationTransformation struct {
	Type       string `json:"type"`
	Field      string `json:"field,omitempty"`
	Expression string `json:"expression,omitempty"`
	MapValue   string `json:"mapValue,omitempty"`
}

// CorrelationCreateRequest is the body for creating a correlation.
type CorrelationCreateRequest struct {
	TargetUID   string            `json:"targetUID"`
	Label       string            `json:"label"`
	Description string            `json:"description,omitempty"`
	Config      CorrelationConfig `json:"config"`
}

// CorrelationUpdateRequest is the body for updating a correlation.
type CorrelationUpdateRequest struct {
	Label       string            `json:"label,omitempty"`
	Description string            `json:"description,omitempty"`
	Config      *CorrelationConfig `json:"config,omitempty"`
}

// CorrelationsResponse is the response from listing correlations.
type CorrelationsResponse struct {
	Correlations []Correlation `json:"correlations,omitempty"`
}

// ListCorrelations returns all correlations.
func (c *Client) ListCorrelations(ctx context.Context) ([]Correlation, error) {
	var result CorrelationsResponse
	resp, err := c.Get(ctx, "/api/datasources/correlations")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return result.Correlations, nil
}

// GetCorrelation returns a correlation by source UID and correlation UID.
func (c *Client) GetCorrelation(ctx context.Context, sourceUID, correlationUID string) (*Correlation, error) {
	var result Correlation
	resp, err := c.Get(ctx, fmt.Sprintf("/api/datasources/uid/%s/correlations/%s", sourceUID, correlationUID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCorrelation creates a new correlation for a datasource.
func (c *Client) CreateCorrelation(ctx context.Context, sourceUID string, req CorrelationCreateRequest) (*Correlation, error) {
	var result Correlation
	resp, err := c.Post(ctx, fmt.Sprintf("/api/datasources/uid/%s/correlations", sourceUID), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCorrelation updates a correlation.
func (c *Client) UpdateCorrelation(ctx context.Context, sourceUID, correlationUID string, req CorrelationUpdateRequest) (*Correlation, error) {
	var result Correlation
	resp, err := c.Patch(ctx, fmt.Sprintf("/api/datasources/uid/%s/correlations/%s", sourceUID, correlationUID), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCorrelation deletes a correlation.
func (c *Client) DeleteCorrelation(ctx context.Context, sourceUID, correlationUID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/datasources/uid/%s/correlations/%s", sourceUID, correlationUID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
