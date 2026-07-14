package client

import (
	"context"
	"fmt"
)

// Folder represents a Grafana folder.
type Folder struct {
	ID        int64  `json:"id"`
	UID       string `json:"uid"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	HasACL    bool   `json:"hasAcl"`
	CanSave   bool   `json:"canSave"`
	CanEdit   bool   `json:"canEdit"`
	CanAdmin  bool   `json:"canAdmin"`
	CanDelete bool   `json:"canDelete"`
	CreatedBy string `json:"createdBy"`
	Created   string `json:"created"`
	UpdatedBy string `json:"updatedBy"`
	Updated   string `json:"updated"`
	Version   int    `json:"version"`
}

// FolderCreateRequest is the body for creating a folder.
type FolderCreateRequest struct {
	UID   string `json:"uid,omitempty"`
	Title string `json:"title"`
}

// FolderUpdateRequest is the body for updating a folder.
type FolderUpdateRequest struct {
	Title     string `json:"title"`
	Version   int    `json:"version,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
}

// FolderPermission represents a folder permission entry.
type FolderPermission struct {
	ID             int64  `json:"id"`
	FolderUID      string `json:"uid"`
	UserID         int64  `json:"userId"`
	UserLogin      string `json:"userLogin"`
	UserEmail      string `json:"userEmail"`
	TeamID         int64  `json:"teamId"`
	Team           string `json:"team"`
	Role           string `json:"role"`
	Permission     int    `json:"permission"`
	PermissionName string `json:"permissionName"`
}

// FolderPermissionUpdate represents a permission update item.
type FolderPermissionUpdate struct {
	UserID     int64  `json:"userId,omitempty"`
	TeamID     int64  `json:"teamId,omitempty"`
	Role       string `json:"role,omitempty"`
	Permission int    `json:"permission"`
}

// FolderPermissionsUpdateRequest is the request to update folder permissions.
type FolderPermissionsUpdateRequest struct {
	Items []FolderPermissionUpdate `json:"items"`
}

// ListFolders returns all folders.
func (c *Client) ListFolders(ctx context.Context, page PageParams) ([]Folder, error) {
	path := page.AppendToPath("/api/folders")
	var results []Folder
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetFolder returns a folder by UID.
func (c *Client) GetFolder(ctx context.Context, uid string) (*Folder, error) {
	var result Folder
	resp, err := c.Get(ctx, "/api/folders/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateFolder creates a new folder.
func (c *Client) CreateFolder(ctx context.Context, req FolderCreateRequest) (*Folder, error) {
	var result Folder
	resp, err := c.Post(ctx, "/api/folders", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateFolder updates an existing folder.
func (c *Client) UpdateFolder(ctx context.Context, uid string, req FolderUpdateRequest) (*Folder, error) {
	var result Folder
	resp, err := c.Put(ctx, "/api/folders/"+uid, req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFolder deletes a folder by UID.
func (c *Client) DeleteFolder(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/folders/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetFolderPermissions returns the permissions for a folder.
func (c *Client) GetFolderPermissions(ctx context.Context, uid string) ([]FolderPermission, error) {
	var results []FolderPermission
	resp, err := c.Get(ctx, fmt.Sprintf("/api/folders/%s/permissions", uid))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateFolderPermissions updates the permissions for a folder.
func (c *Client) UpdateFolderPermissions(ctx context.Context, uid string, req FolderPermissionsUpdateRequest) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/api/folders/%s/permissions", uid), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
