package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	Username   string
	Password   string
	OrgID      int64
	UserAgent  string
}

func NewClient(rc *config.ResolvedConfig) (*Client, error) {
	if rc.URL == "" {
		return nil, fmt.Errorf("grafana URL is required (use --url, GRAFANA_URL, or configure a profile)")
	}

	baseURL := strings.TrimRight(rc.URL, "/")

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   true,
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		Token:     rc.Token,
		Username:  rc.Username,
		Password:  rc.Password,
		OrgID:     rc.OrgID,
		UserAgent: "grafana-cli/" + build.Version,
	}, nil
}

func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, body)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, body)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPatch, path, body)
}

func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, nil)
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	if c.OrgID > 0 {
		req.Header.Set("X-Grafana-Org-Id", fmt.Sprintf("%d", c.OrgID))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	return &Response{HTTPResponse: resp, RequestURL: url}, nil
}
