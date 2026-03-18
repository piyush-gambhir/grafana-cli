package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardVersions(f *cmdutil.Factory) *cobra.Command {
	var (
		limit int
		start int
	)

	cmd := &cobra.Command{
		Use:   "versions <uid>",
		Short: "List dashboard versions",
		Long: `List all versions of a dashboard for version history.

The output includes Version number, Created By, Created date, and the
commit Message. Use this to inspect change history before restoring to
a previous version with "grafana dashboard restore".

Examples:
  # List all versions
  grafana dashboard versions abc123

  # Paginate results
  grafana dashboard versions abc123 --limit 10 --start 5

  # Output as JSON
  grafana dashboard versions abc123 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			versions, err := c.GetDashboardVersions(context.Background(), args[0], limit, start)
			if err != nil {
				return err
			}

			if len(versions) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No versions found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, versions, &output.TableDef{
				Headers: []string{"Version", "Created By", "Created", "Message"},
				RowFunc: func(item interface{}) []string {
					v := item.(client.DashboardVersion)
					return []string{
						fmt.Sprintf("%d", v.Version),
						v.CreatedBy,
						v.Created,
						v.Message,
					}
				},
			})
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of versions to return")
	cmd.Flags().IntVar(&start, "start", 0, "Version ID to start from (for pagination)")

	return cmd
}
