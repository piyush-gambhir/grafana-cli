package client

import (
	"context"
	"fmt"
	"net/url"
)

// DashboardSearchResult represents a single search result.
type DashboardSearchResult struct {
	ID          int64    `json:"id"`
	UID         string   `json:"uid"`
	Title       string   `json:"title"`
	URI         string   `json:"uri"`
	URL         string   `json:"url"`
	Slug        string   `json:"slug"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	IsStarred   bool     `json:"isStarred"`
	FolderID    int64    `json:"folderId"`
	FolderUID   string   `json:"folderUid"`
	FolderTitle string   `json:"folderTitle"`
	FolderURL   string   `json:"folderUrl"`
}

// DashboardFullResponse is the response from GET /api/dashboards/uid/:uid.
type DashboardFullResponse struct {
	Meta      DashboardMeta          `json:"meta"`
	Dashboard map[string]interface{} `json:"dashboard"`
}

// DashboardMeta holds dashboard metadata.
type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	FolderID  int64  `json:"folderId"`
	FolderUID string `json:"folderUid"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	Version   int    `json:"version"`
}

// DashboardCreateRequest is the request body for creating/updating a dashboard.
type DashboardCreateRequest struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	FolderID  int64                  `json:"folderId,omitempty"`
	FolderUID string                 `json:"folderUid,omitempty"`
	Overwrite bool                   `json:"overwrite,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// DashboardCreateResponse is the response from creating/updating a dashboard.
type DashboardCreateResponse struct {
	ID      int64  `json:"id"`
	UID     string `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Version int    `json:"version"`
	Slug    string `json:"slug"`
}

// DashboardVersion represents a dashboard version.
type DashboardVersion struct {
	ID            int    `json:"id"`
	DashboardID   int64  `json:"dashboardId"`
	Version       int    `json:"version"`
	ParentVersion int    `json:"parentVersion"`
	CreatedBy     string `json:"createdBy"`
	Created       string `json:"created"`
	Message       string `json:"message"`
}

// DashboardRestoreRequest is the request to restore a dashboard version.
type DashboardRestoreRequest struct {
	Version int `json:"version"`
}

// DashboardPermission represents a permission entry.
type DashboardPermission struct {
	ID             int64  `json:"id"`
	DashboardID    int64  `json:"dashboardId"`
	DashboardUID   string `json:"dashboardUid"`
	UserID         int64  `json:"userId"`
	UserLogin      string `json:"userLogin"`
	UserEmail      string `json:"userEmail"`
	TeamID         int64  `json:"teamId"`
	Team           string `json:"team"`
	Role           string `json:"role"`
	Permission     int    `json:"permission"`
	PermissionName string `json:"permissionName"`
}

// DashboardPermissionUpdate represents a permission update item.
type DashboardPermissionUpdate struct {
	UserID     int64  `json:"userId,omitempty"`
	TeamID     int64  `json:"teamId,omitempty"`
	Role       string `json:"role,omitempty"`
	Permission int    `json:"permission"`
}

// DashboardPermissionsUpdateRequest is the request to update permissions.
type DashboardPermissionsUpdateRequest struct {
	Items []DashboardPermissionUpdate `json:"items"`
}

// SearchDashboards searches for dashboards.
func (c *Client) SearchDashboards(ctx context.Context, query, tag, folderUID string, page PageParams) ([]DashboardSearchResult, error) {
	v := url.Values{}
	v.Set("type", "dash-db")
	if query != "" {
		v.Set("query", query)
	}
	if tag != "" {
		v.Set("tag", tag)
	}
	if folderUID != "" {
		v.Set("folderUIDs", folderUID)
	}
	page.Apply(v)

	var results []DashboardSearchResult
	resp, err := c.Get(ctx, "/api/search?"+v.Encode())
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetDashboardByUID retrieves a dashboard by UID.
func (c *Client) GetDashboardByUID(ctx context.Context, uid string) (*DashboardFullResponse, error) {
	var result DashboardFullResponse
	resp, err := c.Get(ctx, "/api/dashboards/uid/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDashboard creates or updates a dashboard.
func (c *Client) CreateDashboard(ctx context.Context, req DashboardCreateRequest) (*DashboardCreateResponse, error) {
	var result DashboardCreateResponse
	resp, err := c.Post(ctx, "/api/dashboards/db", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDashboard deletes a dashboard by UID.
func (c *Client) DeleteDashboard(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/dashboards/uid/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetDashboardVersions lists versions of a dashboard.
func (c *Client) GetDashboardVersions(ctx context.Context, dashboardID int64, page PageParams) ([]DashboardVersion, error) {
	path := fmt.Sprintf("/api/dashboards/id/%d/versions", dashboardID)
	path = page.AppendToPath(path)

	var results []DashboardVersion
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// RestoreDashboardVersion restores a dashboard to a specific version.
func (c *Client) RestoreDashboardVersion(ctx context.Context, dashboardID int64, version int) error {
	path := fmt.Sprintf("/api/dashboards/id/%d/restore", dashboardID)
	req := DashboardRestoreRequest{Version: version}
	resp, err := c.Post(ctx, path, req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetDashboardPermissions gets the permissions for a dashboard.
func (c *Client) GetDashboardPermissions(ctx context.Context, uid string) ([]DashboardPermission, error) {
	var results []DashboardPermission
	resp, err := c.Get(ctx, "/api/dashboards/uid/"+uid+"/permissions")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateDashboardPermissions updates the permissions for a dashboard.
func (c *Client) UpdateDashboardPermissions(ctx context.Context, uid string, req DashboardPermissionsUpdateRequest) error {
	resp, err := c.Post(ctx, "/api/dashboards/uid/"+uid+"/permissions", req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
