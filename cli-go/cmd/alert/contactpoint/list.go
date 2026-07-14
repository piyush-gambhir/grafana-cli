package contactpoint

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/output"
)

func newCmdContactPointList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List contact points",
		Long: `List all configured contact points.

The output includes UID, Name, Type (email, slack, etc.), and Provenance.

Examples:
  # List all contact points
  grafana alert contact-point list

  # Output as JSON
  grafana alert contact-point list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListContactPoints(cmd.Context())
			if err != nil {
				return err
			}

			if len(results) == 0 && f.Resolved.Output == "table" {
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
