package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// --- Loki response types ---
//
// Loki query responses come in three shapes depending on the query:
//
//   resultType=streams  → raw log lines    → LokiStream
//   resultType=vector   → instant samples  → PrometheusResult (Prometheus-compatible)
//   resultType=matrix   → range samples    → PrometheusResult (Prometheus-compatible)
//
// A LogQL aggregation like `sum(count_over_time({job="x"}[5m]))` returns
// vector or matrix shapes — not streams. Decoding the whole payload into
// []LokiStream fails ("cannot unmarshal number into Go struct field … of
// type string") and historically dropped data silently on instant queries.
// We deserialize Result as a RawMessage and decode based on ResultType.

// LokiQueryResponse is the top-level response from Loki query endpoints.
type LokiQueryResponse struct {
	Status string        `json:"status"`
	Data   LokiQueryData `json:"data"`
}

// LokiQueryData holds the result type and the raw result payload. Use
// Streams() / Vector() / Matrix() to decode based on ResultType.
type LokiQueryData struct {
	ResultType string          `json:"resultType"`
	Result     json.RawMessage `json:"result"`
}

// Streams decodes the result as []LokiStream when ResultType == "streams".
func (d *LokiQueryData) Streams() ([]LokiStream, error) {
	if d.ResultType != "streams" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "streams")
	}
	var out []LokiStream
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding streams: %w", err)
	}
	return out, nil
}

// Vector decodes the result as []PrometheusResult (instant samples) when
// ResultType == "vector".
func (d *LokiQueryData) Vector() ([]PrometheusResult, error) {
	if d.ResultType != "vector" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "vector")
	}
	var out []PrometheusResult
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding vector: %w", err)
	}
	return out, nil
}

// Matrix decodes the result as []PrometheusResult (range samples) when
// ResultType == "matrix".
func (d *LokiQueryData) Matrix() ([]PrometheusResult, error) {
	if d.ResultType != "matrix" {
		return nil, fmt.Errorf("resultType is %q, expected %q", d.ResultType, "matrix")
	}
	var out []PrometheusResult
	if err := json.Unmarshal(d.Result, &out); err != nil {
		return nil, fmt.Errorf("decoding matrix: %w", err)
	}
	return out, nil
}

// LokiStream represents a single log stream with its entries.
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"` // Each: [timestamp_ns, log_line]
}

// --- Prometheus response types ---

// PrometheusQueryResponse is the top-level response from Prometheus query endpoints.
type PrometheusQueryResponse struct {
	Status string              `json:"status"`
	Data   PrometheusQueryData `json:"data"`
}

// PrometheusQueryData holds the result type and results from a Prometheus query.
type PrometheusQueryData struct {
	ResultType string             `json:"resultType"`
	Result     []PrometheusResult `json:"result"`
}

// PrometheusResult represents a single result from a Prometheus query.
type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value,omitempty"`  // instant: [timestamp, "value"]
	Values [][]interface{}   `json:"values,omitempty"` // range: [[timestamp, "value"], ...]
}

// DatasourceProxyQueryLoki queries a Loki datasource through Grafana's datasource proxy.
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
	default: // range
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

// DatasourceProxyQueryPrometheus queries a Prometheus datasource through Grafana's datasource proxy.
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
	default: // range
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
