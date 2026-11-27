package datasource

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
)

func TestDatasourceQuery_PrometheusRangeDefaultsStep(t *testing.T) {
	var step string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/datasources/uid/prom":
			json.NewEncoder(w).Encode(client.Datasource{UID: "prom", Type: "prometheus"})
		case "/api/datasources/proxy/uid/prom/api/v1/query_range":
			step = r.URL.Query().Get("step")
			w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	var out bytes.Buffer
	f := &cmdutil.Factory{
		IOStreams: cmdutil.IOStreams{Out: &out, ErrOut: &out},
		Resolved:  &config.ResolvedConfig{Output: "table"},
		Client: func() (*client.Client, error) {
			return &client.Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}, nil
		},
	}
	cmd := newCmdDatasourceQuery(f)
	cmd.SetArgs([]string{"prom", "--expr", "up", "--last", "1h"})
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step != "14" {
		t.Errorf("step = %q, want 14 (1h / 250, in seconds)", step)
	}
}

func TestResolveTimeRange_LastEndsAtTo(t *testing.T) {
	start, end, err := resolveTimeRange("30m", "", "2026-11-01T12:00:00Z", "prometheus")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if start != "1793532600.000" || end != "1793534400.000" {
		t.Errorf("got start=%s end=%s, want 30m ending at --to", start, end)
	}
}

func TestParseTime_EpochUnits(t *testing.T) {
	want := time.UnixMilli(1730462400000)
	for name, in := range map[string]string{
		"seconds":      "1730462400",
		"milliseconds": "1730462400000",
		"microseconds": "1730462400000000",
		"nanoseconds":  "1730462400000000000",
	} {
		got, err := parseTime(in)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", name, err)
		}
		if !got.Equal(want) {
			t.Errorf("parseTime(%s) = %v, want %v", name, got, want)
		}
	}
}
