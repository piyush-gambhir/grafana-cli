package client

import (
	"context"
)

type AdminSettings map[string]map[string]string

type AdminStats struct {
	Orgs               int `json:"orgs"`
	Dashboards         int `json:"dashboards"`
	Snapshots          int `json:"snapshots"`
	Tags               int `json:"tags"`
	Datasources        int `json:"datasources"`
	Playlists          int `json:"playlists"`
	Stars              int `json:"stars"`
	Alerts             int `json:"alerts"`
	Users              int `json:"users"`
	Admins             int `json:"admins"`
	Editors            int `json:"editors"`
	Viewers            int `json:"viewers"`
	ActiveUsers        int `json:"activeUsers"`
	ActiveAdmins       int `json:"activeAdmins"`
	ActiveEditors      int `json:"activeEditors"`
	ActiveViewers      int `json:"activeViewers"`
	ActiveSessions     int `json:"activeSessions"`
	DailyActiveUsers   int `json:"dailyActiveUsers"`
	MonthlyActiveUsers int `json:"monthlyActiveUsers"`
}

func (c *Client) GetAdminSettings(ctx context.Context) (AdminSettings, error) {
	var result AdminSettings
	resp, err := c.Get(ctx, "/api/admin/settings")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetAdminStats(ctx context.Context) (*AdminStats, error) {
	var result AdminStats
	resp, err := c.Get(ctx, "/api/admin/stats")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ReloadDashboards(ctx context.Context) error {
	resp, err := c.Post(ctx, "/api/admin/provisioning/dashboards/reload", nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ReloadDatasources(ctx context.Context) error {
	resp, err := c.Post(ctx, "/api/admin/provisioning/datasources/reload", nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ReloadPlugins(ctx context.Context) error {
	resp, err := c.Post(ctx, "/api/admin/provisioning/plugins/reload", nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ReloadAccessControl(ctx context.Context) error {
	resp, err := c.Post(ctx, "/api/admin/provisioning/access-control/reload", nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ReloadAlerting(ctx context.Context) error {
	resp, err := c.Post(ctx, "/api/admin/provisioning/alerting/reload", nil)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
