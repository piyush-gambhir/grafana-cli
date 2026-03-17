package client

import (
	"context"
	"fmt"
	"net/url"
)

// LibraryElement represents a Grafana library element (panel or variable).
type LibraryElement struct {
	ID          int64                  `json:"id"`
	OrgID       int64                  `json:"orgId"`
	FolderID    int64                  `json:"folderId"`
	FolderUID   string                 `json:"folderUid,omitempty"`
	UID         string                 `json:"uid"`
	Name        string                 `json:"name"`
	Kind        int                    `json:"kind"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Model       map[string]interface{} `json:"model"`
	Version     int64                  `json:"version"`
	Meta        LibraryElementMeta     `json:"meta"`
}

// LibraryElementMeta holds library element metadata.
type LibraryElementMeta struct {
	FolderName          string `json:"folderName"`
	FolderUID           string `json:"folderUid"`
	ConnectedDashboards int    `json:"connectedDashboards"`
	Created             string `json:"created"`
	Updated             string `json:"updated"`
	CreatedBy           LibraryElementUser `json:"createdBy"`
	UpdatedBy           LibraryElementUser `json:"updatedBy"`
}

// LibraryElementUser represents a user in library element metadata.
type LibraryElementUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
}

// LibraryElementSearchResult is the result of searching library elements.
type LibraryElementSearchResult struct {
	Result LibraryElementSearchData `json:"result"`
}

// LibraryElementSearchData holds the search result data.
type LibraryElementSearchData struct {
	TotalCount int              `json:"totalCount"`
	Elements   []LibraryElement `json:"elements"`
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
}

// LibraryElementResult wraps a single library element response.
type LibraryElementResult struct {
	Result LibraryElement `json:"result"`
}

// LibraryElementCreateRequest is the body for creating a library element.
type LibraryElementCreateRequest struct {
	FolderID int64                  `json:"folderId,omitempty"`
	FolderUID string                `json:"folderUid,omitempty"`
	Name     string                 `json:"name"`
	Kind     int                    `json:"kind"`
	Model    map[string]interface{} `json:"model"`
}

// LibraryElementUpdateRequest is the body for updating a library element.
type LibraryElementUpdateRequest struct {
	FolderID int64                  `json:"folderId,omitempty"`
	FolderUID string                `json:"folderUid,omitempty"`
	Name     string                 `json:"name"`
	Kind     int                    `json:"kind"`
	Model    map[string]interface{} `json:"model"`
	Version  int64                  `json:"version"`
}

// LibraryElementConnection represents a connection to a dashboard.
type LibraryElementConnection struct {
	ID                int64  `json:"id"`
	Kind              int    `json:"kind"`
	ElementID         int64  `json:"elementId"`
	ConnectionID      int64  `json:"connectionId"`
	ConnectionUID     string `json:"connectionUid,omitempty"`
	Created           string `json:"created"`
	CreatedBy         LibraryElementUser `json:"createdBy"`
}

// LibraryElementConnectionsResult wraps connection results.
type LibraryElementConnectionsResult struct {
	Result []LibraryElementConnection `json:"result"`
}

// ListLibraryElements searches for library elements.
func (c *Client) ListLibraryElements(ctx context.Context, searchString string, kind int, page PageParams) (*LibraryElementSearchResult, error) {
	v := url.Values{}
	if searchString != "" {
		v.Set("searchString", searchString)
	}
	if kind > 0 {
		v.Set("kind", fmt.Sprintf("%d", kind))
	}
	page.Apply(v)

	path := "/api/library-elements?" + v.Encode()
	var result LibraryElementSearchResult
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLibraryElement returns a library element by UID.
func (c *Client) GetLibraryElement(ctx context.Context, uid string) (*LibraryElementResult, error) {
	var result LibraryElementResult
	resp, err := c.Get(ctx, "/api/library-elements/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateLibraryElement creates a new library element.
func (c *Client) CreateLibraryElement(ctx context.Context, req LibraryElementCreateRequest) (*LibraryElementResult, error) {
	var result LibraryElementResult
	resp, err := c.Post(ctx, "/api/library-elements", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateLibraryElement updates a library element.
func (c *Client) UpdateLibraryElement(ctx context.Context, uid string, req LibraryElementUpdateRequest) (*LibraryElementResult, error) {
	var result LibraryElementResult
	resp, err := c.Patch(ctx, "/api/library-elements/"+uid, req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteLibraryElement deletes a library element by UID.
func (c *Client) DeleteLibraryElement(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/library-elements/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetLibraryElementConnections returns dashboards connected to a library element.
func (c *Client) GetLibraryElementConnections(ctx context.Context, uid string) (*LibraryElementConnectionsResult, error) {
	var result LibraryElementConnectionsResult
	resp, err := c.Get(ctx, "/api/library-elements/"+uid+"/connections/")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
