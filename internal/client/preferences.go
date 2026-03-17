package client

import "context"

// Preferences represents user/org preferences.
type Preferences struct {
	Theme           string `json:"theme"`
	HomeDashboardID int64  `json:"homeDashboardId"`
	HomeDashboardUID string `json:"homeDashboardUID,omitempty"`
	Timezone        string `json:"timezone"`
	WeekStart       string `json:"weekStart"`
	Language        string `json:"language,omitempty"`
	QueryHistory    QueryHistoryPreference `json:"queryHistory,omitempty"`
}

// QueryHistoryPreference holds the query history preference.
type QueryHistoryPreference struct {
	HomeTab string `json:"homeTab"`
}

// PreferencesUpdateRequest is the request body for updating preferences.
type PreferencesUpdateRequest struct {
	Theme            string `json:"theme,omitempty"`
	HomeDashboardID  int64  `json:"homeDashboardId,omitempty"`
	HomeDashboardUID string `json:"homeDashboardUID,omitempty"`
	Timezone         string `json:"timezone,omitempty"`
	WeekStart        string `json:"weekStart,omitempty"`
	Language         string `json:"language,omitempty"`
}

// GetPreferences returns the current user's preferences.
func (c *Client) GetPreferences(ctx context.Context) (*Preferences, error) {
	var result Preferences
	resp, err := c.Get(ctx, "/api/user/preferences")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePreferences updates the current user's preferences.
func (c *Client) UpdatePreferences(ctx context.Context, req PreferencesUpdateRequest) error {
	resp, err := c.Put(ctx, "/api/user/preferences", req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetOrgPreferences returns the current org's preferences.
func (c *Client) GetOrgPreferences(ctx context.Context) (*Preferences, error) {
	var result Preferences
	resp, err := c.Get(ctx, "/api/org/preferences")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrgPreferences updates the current org's preferences.
func (c *Client) UpdateOrgPreferences(ctx context.Context, req PreferencesUpdateRequest) error {
	resp, err := c.Put(ctx, "/api/org/preferences", req)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
