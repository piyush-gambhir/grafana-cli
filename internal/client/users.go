package client

import (
	"context"
	"fmt"
	"net/url"
)

// User represents a Grafana user.
type User struct {
	ID            int64  `json:"id"`
	Login         string `json:"login"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	IsAdmin       bool   `json:"isGrafanaAdmin"`
	IsDisabled    bool   `json:"isDisabled"`
	LastSeenAt    string `json:"lastSeenAt,omitempty"`
	LastSeenAtAge string `json:"lastSeenAtAge,omitempty"`
	AuthLabels    []string `json:"authLabels,omitempty"`
	AvatarURL     string `json:"avatarUrl,omitempty"`
}

// UserSearchResult represents the result of a user search.
type UserSearchResult struct {
	TotalCount int    `json:"totalCount"`
	Users      []User `json:"users"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
}

// UserUpdateRequest is the request body for updating a user.
type UserUpdateRequest struct {
	Login string `json:"login,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Theme string `json:"theme,omitempty"`
}

// UserOrg represents an org membership for a user.
type UserOrg struct {
	OrgID int64  `json:"orgId"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// UserTeam represents a team membership for a user.
type UserTeam struct {
	ID          int64  `json:"id"`
	OrgID       int64  `json:"orgId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	MemberCount int    `json:"memberCount"`
	Permission  int    `json:"permission"`
}

// StarResponse is the response from starring/unstarring a dashboard.
type StarResponse struct {
	Message string `json:"message"`
}

// ListUsers returns all users with optional search.
func (c *Client) ListUsers(ctx context.Context, query string, page PageParams) (*UserSearchResult, error) {
	v := url.Values{}
	if query != "" {
		v.Set("query", query)
	}
	page.Apply(v)

	path := "/api/users/search?" + v.Encode()
	var result UserSearchResult
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUser returns a user by ID.
func (c *Client) GetUser(ctx context.Context, id int64) (*User, error) {
	var result User
	resp, err := c.Get(ctx, fmt.Sprintf("/api/users/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LookupUser looks up a user by login or email.
func (c *Client) LookupUser(ctx context.Context, loginOrEmail string) (*User, error) {
	v := url.Values{}
	v.Set("loginOrEmail", loginOrEmail)
	var result User
	resp, err := c.Get(ctx, "/api/users/lookup?"+v.Encode())
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateUser updates a user.
func (c *Client) UpdateUser(ctx context.Context, id int64, req UserUpdateRequest) error {
	resp, err := c.Put(ctx, fmt.Sprintf("/api/users/%d", id), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetUserOrgs returns the organizations a user belongs to.
func (c *Client) GetUserOrgs(ctx context.Context, userID int64) ([]UserOrg, error) {
	var results []UserOrg
	resp, err := c.Get(ctx, fmt.Sprintf("/api/users/%d/orgs", userID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetUserTeams returns the teams a user belongs to.
func (c *Client) GetUserTeams(ctx context.Context, userID int64) ([]UserTeam, error) {
	var results []UserTeam
	resp, err := c.Get(ctx, fmt.Sprintf("/api/users/%d/teams", userID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetCurrentUser returns the authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var result User
	resp, err := c.Get(ctx, "/api/user")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// StarDashboard stars a dashboard for the current user.
func (c *Client) StarDashboard(ctx context.Context, dashboardID int64) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/api/user/stars/dashboard/%d", dashboardID), nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// UnstarDashboard removes a star from a dashboard for the current user.
func (c *Client) UnstarDashboard(ctx context.Context, dashboardID int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/user/stars/dashboard/%d", dashboardID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
