package annotation

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

// parseAnnotationTime accepts RFC3339, Unix seconds (int or float),
// epoch milliseconds, or epoch nanoseconds and returns epoch ms. Grafana's
// annotation API expects milliseconds — historically this command took
// an int64 flag which silently treated RFC3339 input as 0 and returned
// annotations from 1970. Accepting both formats removes that footgun.
func parseAnnotationTime(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UnixMilli(), nil
	}
	// Numeric input — decide between seconds, milliseconds, and nanoseconds
	// using magnitude thresholds. Values < 1e12 are seconds (through year
	// 33658), 1e12-1e15 are milliseconds, >= 1e15 are nanoseconds.
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		switch {
		case n >= 1_000_000_000_000_000: // 1e15 → nanoseconds
			return n / 1_000_000, nil
		case n >= 1_000_000_000_000: // 1e12 → milliseconds
			return n, nil
		case n > 0: // seconds
			return n * 1000, nil
		default:
			return 0, nil
		}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f * 1000), nil
	}
	return 0, fmt.Errorf("cannot parse %q as RFC3339 or a numeric epoch", s)
}

func newCmdAnnotationList(f *cmdutil.Factory) *cobra.Command {
	var (
		dashboardID int64
		panelID     int64
		fromStr     string
		toStr       string
		tags        []string
		limit       int64
		annType     string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List annotations",
		Aliases: []string{"ls"},
		Long: `List annotations in the current organization with optional filters.

The output includes ID, Dashboard ID, Text, Tags, and Time (epoch ms).
Multiple filters can be combined to narrow results.

Time values (--from, --to) accept RFC3339 timestamps, Unix seconds,
epoch milliseconds, or epoch nanoseconds. The CLI normalises to the
millisecond epoch that Grafana's API expects. The --type flag accepts
"annotation" or "alert" to filter by annotation source.

Note: annotations are only recorded for alerts managed by Grafana. Alerts
raised by downstream systems (for example native CubeAPM alerting) do not
appear here — check the originating system's alert surface instead.

Examples:
  # List all annotations (default limit 100)
  grafana annotation list

  # List annotations for a specific dashboard
  grafana annotation list --dashboard-id 42

  # List annotations for a specific panel
  grafana annotation list --dashboard-id 42 --panel-id 3

  # List annotations within a time range (RFC3339)
  grafana annotation list --from 2024-01-01T00:00:00Z --to 2024-01-02T00:00:00Z

  # List annotations within a time range (epoch milliseconds)
  grafana annotation list --from 1609459200000 --to 1609545600000

  # Filter by tags
  grafana annotation list --tags deploy,release

  # Filter by annotation type (annotation or alert)
  grafana annotation list --type alert

  # Increase the result limit
  grafana annotation list --limit 500

  # Output as JSON
  grafana annotation list -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := parseAnnotationTime(fromStr)
			if err != nil {
				return fmt.Errorf("--from: %w", err)
			}
			to, err := parseAnnotationTime(toStr)
			if err != nil {
				return fmt.Errorf("--to: %w", err)
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListAnnotations(context.Background(), dashboardID, panelID, from, to, tags, limit, annType)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No annotations found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Dashboard ID", "Text", "Tags", "Time"},
				RowFunc: func(item interface{}) []string {
					a := item.(client.Annotation)
					return []string{
						fmt.Sprintf("%d", a.ID),
						fmt.Sprintf("%d", a.DashboardID),
						a.Text,
						strings.Join(a.Tags, ", "),
						fmt.Sprintf("%d", a.Time),
					}
				},
			})
		},
	}

	cmd.Flags().Int64Var(&dashboardID, "dashboard-id", 0, "Filter by dashboard ID")
	cmd.Flags().Int64Var(&panelID, "panel-id", 0, "Filter by panel ID")
	cmd.Flags().StringVar(&fromStr, "from", "", "Start time (RFC3339, Unix seconds, epoch ms, or epoch ns)")
	cmd.Flags().StringVar(&toStr, "to", "", "End time (RFC3339, Unix seconds, epoch ms, or epoch ns)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Filter by tags (comma-separated)")
	cmd.Flags().Int64Var(&limit, "limit", 100, "Maximum number of annotations to return")
	cmd.Flags().StringVar(&annType, "type", "", "Filter by type: annotation or alert")

	return cmd
}
