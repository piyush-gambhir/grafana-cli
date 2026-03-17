package contactpoint

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdContactPointList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List contact points",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListContactPoints(context.Background())
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No contact points found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"UID", "Name", "Type", "Provenance"},
				RowFunc: func(item interface{}) []string {
					cp := item.(client.ContactPoint)
					return []string{cp.UID, cp.Name, cp.Type, cp.Provenance}
				},
			})
		},
	}
}
