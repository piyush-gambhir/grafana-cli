package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type LokiQueryResponse struct {
	Status string        `json:"status"`
	Data   LokiQueryData `json:"data"`
}

type LokiQueryData struct {
	ResultType string          `json:"resultType"`
	Result     json.RawMessage `json:"result"`
}

func (d *LokiQueryData) Streams() ([]LokiStream, error) {
	if d.ResultType != "streams" && d.ResultType != "" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "streams")
	}
	if len(d.Result) == 0 {
		return nil, nil
	}
	var out []LokiStream
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding streams: %w", err)
	}
	return out, nil
}

func (d *LokiQueryData) Vector() ([]PrometheusResult, error) {
	if d.ResultType != "vector" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "vector")
	}
	if len(d.Result) == 0 {
		return nil, nil
	}
	var out []PrometheusResult
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding vector: %w", err)
	}
	return out, nil
}

func (d *LokiQueryData) Matrix() ([]PrometheusResult, error) {
	if d.ResultType != "matrix" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "matrix")
	}
	if len(d.Result) == 0 {
		return nil, nil
	}
	var out []PrometheusResult
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding matrix: %w", err)
	}
	return out, nil
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type PrometheusQueryResponse struct {
	Status string              `json:"status"`
	Data   PrometheusQueryData `json:"data"`
}

type PrometheusQueryData struct {
	ResultType string             `json:"resultType"`
	Result     []PrometheusResult `json:"result"`
}

func (d *PrometheusQueryData) UnmarshalJSON(b []byte) error {
	var raw struct {
		ResultType string          `json:"resultType"`
		Result     json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	d.ResultType = raw.ResultType
	d.Result = nil
	if len(raw.Result) == 0 || string(raw.Result) == "null" {
		return nil
	}
	switch raw.ResultType {
	case "scalar", "string":
		var pair []interface{}
		if err := json.Unmarshal(raw.Result, &pair); err != nil {
			return fmt.Errorf("decoding %s result: %w", raw.ResultType, err)
		}
		d.Result = []PrometheusResult{{Value: pair}}
		return nil
	default:
		if err := json.Unmarshal(raw.Result, &d.Result); err != nil {
			return fmt.Errorf("decoding %s result: %w", raw.ResultType, err)
		}
		return nil
	}
}

type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value,omitempty"`
	Values [][]interface{}   `json:"values,omitempty"`
}

func (c *Client) DatasourceProxyQueryLoki(ctx context.Context, uid string, queryType string, expr string, start, end string, limit int, direction string) (*LokiQueryResponse, error) {
	v := url.Values{}
	v.Set("query", expr)
	if limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", limit))
	}
	if direction != "" {
		v.Set("direction", direction)
	}

	var endpoint string
	switch queryType {
	case "instant":
		endpoint = "loki/api/v1/query"
		if end != "" {
			v.Set("time", end)
		}
	default:
		endpoint = "loki/api/v1/query_range"
		if start != "" {
			v.Set("start", start)
		}
		if end != "" {
			v.Set("end", end)
		}
	}

	path := fmt.Sprintf("/api/datasources/proxy/uid/%s/%s?%s", uid, endpoint, v.Encode())

	var result LokiQueryResponse
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DatasourceProxyQueryPrometheus(ctx context.Context, uid string, queryType string, expr string, start, end, step string) (*PrometheusQueryResponse, error) {
	v := url.Values{}
	v.Set("query", expr)

	var endpoint string
	switch queryType {
	case "instant":
		endpoint = "api/v1/query"
		if end != "" {
			v.Set("time", end)
		}
	default:
		endpoint = "api/v1/query_range"
		if start != "" {
			v.Set("start", start)
		}
		if end != "" {
			v.Set("end", end)
		}
		if step != "" {
			v.Set("step", step)
		}
	}

	path := fmt.Sprintf("/api/datasources/proxy/uid/%s/%s?%s", uid, endpoint, v.Encode())

	var result PrometheusQueryResponse
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
