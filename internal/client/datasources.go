package client

import (
	"context"
	"fmt"
)

// Datasource represents a Grafana data source.
type Datasource struct {
	ID                int64                  `json:"id"`
	UID               string                 `json:"uid"`
	OrgID             int64                  `json:"orgId"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	TypeName          string                 `json:"typeName,omitempty"`
	TypeLogoURL       string                 `json:"typeLogoUrl,omitempty"`
	Access            string                 `json:"access"`
	URL               string                 `json:"url"`
	User              string                 `json:"user,omitempty"`
	Database          string                 `json:"database,omitempty"`
	BasicAuth         bool                   `json:"basicAuth"`
	BasicAuthUser     string                 `json:"basicAuthUser,omitempty"`
	WithCredentials   bool                   `json:"withCredentials"`
	IsDefault         bool                   `json:"isDefault"`
	ReadOnly          bool                   `json:"readOnly"`
	JSONData          map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONFields  map[string]bool        `json:"secureJsonFields,omitempty"`
}

// DatasourceCreateRequest is the body for creating a data source.
type DatasourceCreateRequest struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Access          string                 `json:"access"`
	URL             string                 `json:"url,omitempty"`
	User            string                 `json:"user,omitempty"`
	Database        string                 `json:"database,omitempty"`
	BasicAuth       bool                   `json:"basicAuth,omitempty"`
	BasicAuthUser   string                 `json:"basicAuthUser,omitempty"`
	WithCredentials bool                   `json:"withCredentials,omitempty"`
	IsDefault       bool                   `json:"isDefault,omitempty"`
	JSONData        map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONData  map[string]string      `json:"secureJsonData,omitempty"`
}

// DatasourceCreateResponse is the response from creating a data source.
type DatasourceCreateResponse struct {
	Datasource Datasource `json:"datasource"`
	ID         int64      `json:"id"`
	Message    string     `json:"message"`
	Name       string     `json:"name"`
}

// ListDatasources returns all data sources.
func (c *Client) ListDatasources(ctx context.Context) ([]Datasource, error) {
	var results []Datasource
	resp, err := c.Get(ctx, "/api/datasources")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetDatasource returns a data source by ID.
func (c *Client) GetDatasource(ctx context.Context, id int64) (*Datasource, error) {
	var result Datasource
	resp, err := c.Get(ctx, fmt.Sprintf("/api/datasources/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDatasourceByUID returns a data source by UID.
func (c *Client) GetDatasourceByUID(ctx context.Context, uid string) (*Datasource, error) {
	var result Datasource
	resp, err := c.Get(ctx, "/api/datasources/uid/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDatasource creates a new data source.
func (c *Client) CreateDatasource(ctx context.Context, req DatasourceCreateRequest) (*DatasourceCreateResponse, error) {
	var result DatasourceCreateResponse
	resp, err := c.Post(ctx, "/api/datasources", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDatasource updates an existing data source.
func (c *Client) UpdateDatasource(ctx context.Context, id int64, req DatasourceCreateRequest) (*DatasourceCreateResponse, error) {
	var result DatasourceCreateResponse
	resp, err := c.Put(ctx, fmt.Sprintf("/api/datasources/%d", id), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDatasource deletes a data source by ID.
func (c *Client) DeleteDatasource(ctx context.Context, id int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/datasources/%d", id))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// DeleteDatasourceByUID deletes a data source by UID.
func (c *Client) DeleteDatasourceByUID(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/datasources/uid/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
