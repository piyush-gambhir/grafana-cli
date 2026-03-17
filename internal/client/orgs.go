package client

import (
	"context"
	"fmt"
)

// Org represents a Grafana organization.
type Org struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// OrgDetail represents the detail of an organization.
type OrgDetail struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address Address `json:"address"`
}

// Address holds org address info.
type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

// OrgCreateRequest is the request body for creating an org.
type OrgCreateRequest struct {
	Name string `json:"name"`
}

// OrgCreateResponse is the response from creating an org.
type OrgCreateResponse struct {
	OrgID   int64  `json:"orgId"`
	Message string `json:"message"`
}

// OrgUpdateRequest is the request body for updating an org.
type OrgUpdateRequest struct {
	Name string `json:"name"`
}

// OrgUser represents a user in an organization.
type OrgUser struct {
	OrgID     int64  `json:"orgId"`
	UserID    int64  `json:"userId"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	LastSeenAt string `json:"lastSeenAt,omitempty"`
}

// OrgUserAddRequest is the request to add a user to an org.
type OrgUserAddRequest struct {
	LoginOrEmail string `json:"loginOrEmail"`
	Role         string `json:"role"`
}

// OrgUserUpdateRequest is the request to update a user's role in an org.
type OrgUserUpdateRequest struct {
	Role string `json:"role"`
}

// ListOrgs returns all organizations.
func (c *Client) ListOrgs(ctx context.Context, page PageParams) ([]Org, error) {
	path := page.AppendToPath("/api/orgs")
	var results []Org
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetOrg returns an organization by ID.
func (c *Client) GetOrg(ctx context.Context, id int64) (*OrgDetail, error) {
	var result OrgDetail
	resp, err := c.Get(ctx, fmt.Sprintf("/api/orgs/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrg creates a new organization.
func (c *Client) CreateOrg(ctx context.Context, req OrgCreateRequest) (*OrgCreateResponse, error) {
	var result OrgCreateResponse
	resp, err := c.Post(ctx, "/api/orgs", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrg updates an organization.
func (c *Client) UpdateOrg(ctx context.Context, id int64, req OrgUpdateRequest) error {
	resp, err := c.Put(ctx, fmt.Sprintf("/api/orgs/%d", id), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// DeleteOrg deletes an organization by ID.
func (c *Client) DeleteOrg(ctx context.Context, id int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/orgs/%d", id))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetCurrentOrg returns the current organization.
func (c *Client) GetCurrentOrg(ctx context.Context) (*OrgDetail, error) {
	var result OrgDetail
	resp, err := c.Get(ctx, "/api/org/")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SwitchOrg switches the user's active organization.
func (c *Client) SwitchOrg(ctx context.Context, orgID int64) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/api/user/using/%d", orgID), nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ListOrgUsers returns users in an organization.
func (c *Client) ListOrgUsers(ctx context.Context, orgID int64) ([]OrgUser, error) {
	var results []OrgUser
	resp, err := c.Get(ctx, fmt.Sprintf("/api/orgs/%d/users", orgID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// AddOrgUser adds a user to an organization.
func (c *Client) AddOrgUser(ctx context.Context, orgID int64, req OrgUserAddRequest) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/api/orgs/%d/users", orgID), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// UpdateOrgUser updates a user's role in an organization.
func (c *Client) UpdateOrgUser(ctx context.Context, orgID, userID int64, req OrgUserUpdateRequest) error {
	resp, err := c.Patch(ctx, fmt.Sprintf("/api/orgs/%d/users/%d", orgID, userID), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// RemoveOrgUser removes a user from an organization.
func (c *Client) RemoveOrgUser(ctx context.Context, orgID, userID int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/orgs/%d/users/%d", orgID, userID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
