package annotation

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAnnotationList(f *cmdutil.Factory) *cobra.Command {
	var (
		dashboardID int64
		panelID     int64
		from        int64
		to          int64
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

Time values (--from, --to) are epoch milliseconds. The --type flag accepts
"annotation" or "alert" to filter by annotation source.

Examples:
  # List all annotations (default limit 100)
  grafana annotation list

  # List annotations for a specific dashboard
  grafana annotation list --dashboard-id 42

  # List annotations for a specific panel
  grafana annotation list --dashboard-id 42 --panel-id 3

  # List annotations within a time range
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
	cmd.Flags().Int64Var(&from, "from", 0, "Start time (epoch ms)")
	cmd.Flags().Int64Var(&to, "to", 0, "End time (epoch ms)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Filter by tags (comma-separated)")
	cmd.Flags().Int64Var(&limit, "limit", 100, "Maximum number of annotations to return")
	cmd.Flags().StringVar(&annType, "type", "", "Filter by type: annotation or alert")

	return cmd
}
