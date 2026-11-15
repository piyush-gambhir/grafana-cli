package client

import (
	"context"
	"fmt"
	"net/url"
)

type Team struct {
	ID          int64  `json:"id"`
	OrgID       int64  `json:"orgId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	MemberCount int    `json:"memberCount"`
	Permission  int    `json:"permission"`
}

type TeamSearchResult struct {
	TotalCount int    `json:"totalCount"`
	Teams      []Team `json:"teams"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
}

type TeamCreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type TeamCreateResponse struct {
	TeamID  int64  `json:"teamId"`
	Message string `json:"message"`
}

type TeamUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type TeamMember struct {
	OrgID      int64    `json:"orgId"`
	TeamID     int64    `json:"teamId"`
	UserID     int64    `json:"userId"`
	Login      string   `json:"login"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	AvatarURL  string   `json:"avatarUrl,omitempty"`
	Labels     []string `json:"labels,omitempty"`
	Permission int      `json:"permission"`
}

type TeamMemberAddRequest struct {
	UserID int64 `json:"userId"`
}

type TeamPreferences struct {
	Theme           string `json:"theme"`
	HomeDashboardID int64  `json:"homeDashboardId"`
	Timezone        string `json:"timezone"`
	WeekStart       string `json:"weekStart"`
}

func (c *Client) ListTeams(ctx context.Context, query string, page PageParams) (*TeamSearchResult, error) {
	v := url.Values{}
	if query != "" {
		v.Set("query", query)
	}
	page.Apply(v)

	path := "/api/teams/search?" + v.Encode()
	var result TeamSearchResult
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetTeam(ctx context.Context, id int64) (*Team, error) {
	var result Team
	resp, err := c.Get(ctx, fmt.Sprintf("/api/teams/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateTeam(ctx context.Context, req TeamCreateRequest) (*TeamCreateResponse, error) {
	var result TeamCreateResponse
	resp, err := c.Post(ctx, "/api/teams", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateTeam(ctx context.Context, id int64, req TeamUpdateRequest) error {
	resp, err := c.Put(ctx, fmt.Sprintf("/api/teams/%d", id), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) DeleteTeam(ctx context.Context, id int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/teams/%d", id))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ListTeamMembers(ctx context.Context, teamID int64) ([]TeamMember, error) {
	var results []TeamMember
	resp, err := c.Get(ctx, fmt.Sprintf("/api/teams/%d/members", teamID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) AddTeamMember(ctx context.Context, teamID int64, req TeamMemberAddRequest) error {
	resp, err := c.Post(ctx, fmt.Sprintf("/api/teams/%d/members", teamID), req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) RemoveTeamMember(ctx context.Context, teamID, userID int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/teams/%d/members/%d", teamID, userID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) GetTeamPreferences(ctx context.Context, teamID int64) (*TeamPreferences, error) {
	var result TeamPreferences
	resp, err := c.Get(ctx, fmt.Sprintf("/api/teams/%d/preferences", teamID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateTeamPreferences(ctx context.Context, teamID int64, prefs TeamPreferences) error {
	resp, err := c.Put(ctx, fmt.Sprintf("/api/teams/%d/preferences", teamID), prefs)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
