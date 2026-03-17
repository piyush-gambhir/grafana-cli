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
	return &cobra.Command{
		Use:     "list",
		Short:   "List snapshots",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListSnapshots(context.Background())
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
}
