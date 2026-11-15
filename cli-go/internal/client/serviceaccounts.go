package client

import (
	"context"
	"fmt"
	"net/url"
)

type ServiceAccount struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Login      string `json:"login"`
	OrgID      int64  `json:"orgId"`
	IsDisabled bool   `json:"isDisabled"`
	Role       string `json:"role"`
	Tokens     int64  `json:"tokens"`
	AvatarURL  string `json:"avatarUrl,omitempty"`
}

type ServiceAccountSearchResult struct {
	TotalCount      int              `json:"totalCount"`
	ServiceAccounts []ServiceAccount `json:"serviceAccounts"`
	Page            int              `json:"page"`
	PerPage         int              `json:"perPage"`
}

type ServiceAccountCreateRequest struct {
	Name       string `json:"name"`
	Role       string `json:"role,omitempty"`
	IsDisabled bool   `json:"isDisabled,omitempty"`
}

type ServiceAccountUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	Role       string `json:"role,omitempty"`
	IsDisabled *bool  `json:"isDisabled,omitempty"`
}

type ServiceAccountToken struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key,omitempty"`
	Created    string `json:"created,omitempty"`
	Expiration string `json:"expiration,omitempty"`
	LastUsedAt string `json:"lastUsedAt,omitempty"`
	HasExpired bool   `json:"hasExpired"`
}

type ServiceAccountTokenCreateRequest struct {
	Name          string `json:"name"`
	SecondsToLive int64  `json:"secondsToLive,omitempty"`
}

func (c *Client) ListServiceAccounts(ctx context.Context, query string, page PageParams) (*ServiceAccountSearchResult, error) {
	v := url.Values{}
	if query != "" {
		v.Set("query", query)
	}
	page.Apply(v)

	path := "/api/serviceaccounts/search?" + v.Encode()
	var result ServiceAccountSearchResult
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetServiceAccount(ctx context.Context, id int64) (*ServiceAccount, error) {
	var result ServiceAccount
	resp, err := c.Get(ctx, fmt.Sprintf("/api/serviceaccounts/%d", id))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateServiceAccount(ctx context.Context, req ServiceAccountCreateRequest) (*ServiceAccount, error) {
	var result ServiceAccount
	resp, err := c.Post(ctx, "/api/serviceaccounts", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateServiceAccount(ctx context.Context, id int64, req ServiceAccountUpdateRequest) (*ServiceAccount, error) {
	var result ServiceAccount
	resp, err := c.Patch(ctx, fmt.Sprintf("/api/serviceaccounts/%d", id), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteServiceAccount(ctx context.Context, id int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/serviceaccounts/%d", id))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ListServiceAccountTokens(ctx context.Context, saID int64) ([]ServiceAccountToken, error) {
	var results []ServiceAccountToken
	resp, err := c.Get(ctx, fmt.Sprintf("/api/serviceaccounts/%d/tokens", saID))
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) CreateServiceAccountToken(ctx context.Context, saID int64, req ServiceAccountTokenCreateRequest) (*ServiceAccountToken, error) {
	var result ServiceAccountToken
	resp, err := c.Post(ctx, fmt.Sprintf("/api/serviceaccounts/%d/tokens", saID), req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteServiceAccountToken(ctx context.Context, saID, tokenID int64) error {
	resp, err := c.Delete(ctx, fmt.Sprintf("/api/serviceaccounts/%d/tokens/%d", saID, tokenID))
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
