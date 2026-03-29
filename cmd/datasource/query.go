package datasource

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

// queryResultRow is a flattened row for table rendering.
type queryResultRow struct {
	Timestamp string
	Labels    string
	Value     string
}

func newCmdDatasourceQuery(f *cmdutil.Factory) *cobra.Command {
	var (
		expr      string
		last      string
		from      string
		to        string
		limit     int
		direction string
		step      string
		queryType string
	)

	cmd := &cobra.Command{
		Use:   "query <uid>",
		Short: "Query a datasource via Grafana proxy",
		Long: `Query a Loki or Prometheus datasource through Grafana's datasource proxy API.

The datasource type is auto-detected from its UID. Supported types: loki, prometheus.

For Loki, the expression should be a LogQL query. For Prometheus, a PromQL query.
Time range defaults to the last 1 hour. Use --last for relative ranges or
--from/--to for absolute times (RFC3339 or Unix epoch).

Examples:
  # Query Loki logs
  grafana datasource query bdzrc5mgatywwe --expr '{app="media-service"}'

  # Search for errors in the last 30 minutes
  grafana datasource query bdzrc5mgatywwe --expr '{app="api"} |= "error"' --last 30m

  # Query Prometheus metrics
  grafana datasource query adzsonhnrp8u8f --expr 'up' --query-type instant

  # Prometheus range query with step
  grafana datasource query adzsonhnrp8u8f --expr 'rate(http_requests_total[5m])' --last 1h --step 30s

  # Output as JSON
  grafana datasource query bdzrc5mgatywwe --expr '{app="api"}' -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if expr == "" {
				return fmt.Errorf("--expr is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			ctx := context.Background()
			uid := args[0]

			// Auto-detect datasource type.
			ds, err := c.GetDatasourceByUID(ctx, uid)
			if err != nil {
				return fmt.Errorf("fetching datasource: %w", err)
			}

			dsType := ds.Type
			if dsType != "loki" && dsType != "prometheus" {
				return fmt.Errorf("datasource %q is type %q; query supports loki and prometheus", uid, dsType)
			}

			// Resolve time range.
			startStr, endStr, err := resolveTimeRange(last, from, to, dsType)
			if err != nil {
				return err
			}

			switch dsType {
			case "loki":
				return queryLoki(ctx, c, f, uid, queryType, expr, startStr, endStr, limit, direction)
			case "prometheus":
				return queryPrometheus(ctx, c, f, uid, queryType, expr, startStr, endStr, step)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&expr, "expr", "e", "", "Query expression (LogQL or PromQL)")
	cmd.Flags().StringVar(&last, "last", "1h", "Lookback duration (e.g. 1h, 30m)")
	cmd.Flags().StringVar(&from, "from", "", "Start time (RFC3339 or Unix epoch)")
	cmd.Flags().StringVar(&to, "to", "", "End time (RFC3339 or Unix epoch)")
	cmd.Flags().IntVar(&limit, "limit", 100, "Max entries to return (Loki only)")
	cmd.Flags().StringVar(&direction, "direction", "backward", "Log ordering: backward or forward (Loki only)")
	cmd.Flags().StringVar(&step, "step", "", "Query resolution step (e.g. 15s, 1m; Prometheus only)")
	cmd.Flags().StringVar(&queryType, "query-type", "range", "Query type: range or instant")

	return cmd
}

func queryLoki(ctx context.Context, c *client.Client, f *cmdutil.Factory, uid, queryType, expr, start, end string, limit int, direction string) error {
	result, err := c.DatasourceProxyQueryLoki(ctx, uid, queryType, expr, start, end, limit, direction)
	if err != nil {
		return err
	}

	if result.Status != "success" {
		return fmt.Errorf("loki query failed with status: %s", result.Status)
	}

	// For JSON/YAML, output the raw response.
	if f.Resolved.Output != "table" && f.Resolved.Output != "" {
		return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
	}

	// Flatten streams into rows for table output.
	var rows []queryResultRow
	for _, stream := range result.Data.Result {
		labels := formatLabels(stream.Stream)
		for _, entry := range stream.Values {
			if len(entry) < 2 {
				continue
			}
			ts := formatNanosTimestamp(entry[0])
			rows = append(rows, queryResultRow{
				Timestamp: ts,
				Labels:    labels,
				Value:     entry[1],
			})
		}
	}

	if len(rows) == 0 {
		fmt.Fprintln(f.IOStreams.Out, "No results found.")
		return nil
	}

	return output.Print(f.IOStreams.Out, "table", rows, &output.TableDef{
		Headers: []string{"Timestamp", "Labels", "Line"},
		RowFunc: func(item interface{}) []string {
			r := item.(queryResultRow)
			line := r.Value
			if len(line) > 200 {
				line = line[:200] + "..."
			}
			return []string{r.Timestamp, r.Labels, line}
		},
	})
}

func queryPrometheus(ctx context.Context, c *client.Client, f *cmdutil.Factory, uid, queryType, expr, start, end, step string) error {
	result, err := c.DatasourceProxyQueryPrometheus(ctx, uid, queryType, expr, start, end, step)
	if err != nil {
		return err
	}

	if result.Status != "success" {
		return fmt.Errorf("prometheus query failed with status: %s", result.Status)
	}

	// For JSON/YAML, output the raw response.
	if f.Resolved.Output != "table" && f.Resolved.Output != "" {
		return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
	}

	// Flatten results into rows.
	var rows []queryResultRow
	for _, r := range result.Data.Result {
		labels := formatLabels(r.Metric)

		if r.Value != nil && len(r.Value) >= 2 {
			// Instant query result.
			ts := formatPromTimestamp(r.Value[0])
			val := fmt.Sprintf("%v", r.Value[1])
			rows = append(rows, queryResultRow{Timestamp: ts, Labels: labels, Value: val})
		}

		for _, v := range r.Values {
			if len(v) < 2 {
				continue
			}
			ts := formatPromTimestamp(v[0])
			val := fmt.Sprintf("%v", v[1])
			rows = append(rows, queryResultRow{Timestamp: ts, Labels: labels, Value: val})
		}
	}

	if len(rows) == 0 {
		fmt.Fprintln(f.IOStreams.Out, "No results found.")
		return nil
	}

	return output.Print(f.IOStreams.Out, "table", rows, &output.TableDef{
		Headers: []string{"Metric", "Timestamp", "Value"},
		RowFunc: func(item interface{}) []string {
			r := item.(queryResultRow)
			return []string{r.Labels, r.Timestamp, r.Value}
		},
	})
}

// resolveTimeRange computes start/end strings from flags.
// For Loki, returns nanosecond epoch strings. For Prometheus, returns Unix second strings.
func resolveTimeRange(last, from, to, dsType string) (string, string, error) {
	now := time.Now()

	if from != "" {
		startTime, err := parseTime(from)
		if err != nil {
			return "", "", fmt.Errorf("invalid --from: %w", err)
		}
		endTime := now
		if to != "" {
			endTime, err = parseTime(to)
			if err != nil {
				return "", "", fmt.Errorf("invalid --to: %w", err)
			}
		}
		return formatTime(startTime, dsType), formatTime(endTime, dsType), nil
	}

	dur, err := time.ParseDuration(last)
	if err != nil {
		return "", "", fmt.Errorf("invalid --last duration: %w", err)
	}
	startTime := now.Add(-dur)
	return formatTime(startTime, dsType), formatTime(now, dsType), nil
}

func parseTime(s string) (time.Time, error) {
	// Try RFC3339 first.
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try Unix epoch (seconds as float or integer).
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		sec := int64(f)
		nsec := int64((f - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	}
	// Try nanosecond epoch.
	if ns, err := strconv.ParseInt(s, 10, 64); err == nil && ns > 1e15 {
		return time.Unix(0, ns), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as RFC3339 or Unix epoch", s)
}

func formatTime(t time.Time, dsType string) string {
	switch dsType {
	case "loki":
		return fmt.Sprintf("%d", t.UnixNano())
	default: // prometheus
		return fmt.Sprintf("%.3f", float64(t.UnixNano())/1e9)
	}
}

func formatNanosTimestamp(ns string) string {
	n, err := strconv.ParseInt(ns, 10, 64)
	if err != nil {
		return ns
	}
	return time.Unix(0, n).UTC().Format("2006-01-02T15:04:05.000Z")
}

func formatPromTimestamp(v interface{}) string {
	switch ts := v.(type) {
	case float64:
		sec := int64(ts)
		nsec := int64((ts - float64(sec)) * 1e9)
		return time.Unix(sec, nsec).UTC().Format("2006-01-02T15:04:05Z")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatLabels(m map[string]string) string {
	if len(m) == 0 {
		return "{}"
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s=%q", k, v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
