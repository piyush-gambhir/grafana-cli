package snapshot

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdSnapshotList(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List snapshots",
		Aliases: []string{"ls"},
		Long: `List dashboard snapshots for the current organization.

The output includes ID, Name, Key, whether the snapshot is external, and
its expiration date. The --limit flag controls how many snapshots to return
(the Grafana API defaults to returning all snapshots).

The snapshot key is used to retrieve or delete a specific snapshot.

Examples:
  # List all snapshots
  grafana snapshot list

  # List the first 10 snapshots
  grafana snapshot list --limit 10

  # Output as JSON
  grafana snapshot list -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListSnapshots(context.Background(), limit)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No snapshots found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Name", "Key", "External", "Expires"},
				RowFunc: func(item interface{}) []string {
					s := item.(client.Snapshot)
					return []string{
						fmt.Sprintf("%d", s.ID),
						s.Name,
						s.Key,
						fmt.Sprintf("%v", s.External),
						s.Expires,
					}
				},
			})
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of snapshots to return (0 = all)")

	return cmd
}
