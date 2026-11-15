package client

import (
	"context"
	"fmt"
)

type Correlation struct {
	UID         string            `json:"uid"`
	SourceUID   string            `json:"sourceUID"`
	TargetUID   string            `json:"targetUID"`
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Config      CorrelationConfig `json:"config"`
	Provisioned bool              `json:"provisioned,omitempty"`
}

type CorrelationConfig struct {
	Type            string                      `json:"type"`
	Field           string                      `json:"field"`
	Target          map[string]interface{}      `json:"target,omitempty"`
	Transformations []CorrelationTransformation `json:"transformations,omitempty"`
}

type CorrelationTransformation struct {
	Type       string `json:"type"`
	Field      string `json:"field,omitempty"`
	Expression string `json:"expression,omitempty"`
	MapValue   string `json:"mapValue,omitempty"`
}

type CorrelationCreateRequest struct {
	TargetUID   string            `json:"targetUID"`
	Label       string            `json:"label"`
	Description string            `json:"description,omitempty"`
	Config      CorrelationConfig `json:"config"`
}

type CorrelationUpdateRequest struct {
	Label       string             `json:"label,omitempty"`
	Description string             `json:"description,omitempty"`
	Config      *CorrelationConfig `json:"config,omitempty"`
}

type CorrelationsResponse struct {
	Correlations []Correlation `json:"correlations,omitempty"`
}

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

func (c *Client) DeleteCorrelation(ctx context.Context, sourceUID, correlationUID string) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/datasources/uid/%s/correlations/%s", sourceUID, correlationUID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
