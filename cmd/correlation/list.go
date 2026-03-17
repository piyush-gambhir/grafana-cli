package correlation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdCorrelationList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List correlations",
		Long: `List all correlations.

The output includes UID, Source UID, Target UID, Label, and Description.

Examples:
  # List all correlations
  grafana correlation list

  # Output as JSON
  grafana correlation list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListCorrelations(context.Background())
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No correlations found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"UID", "Source UID", "Target UID", "Label", "Description"},
				RowFunc: func(item interface{}) []string {
					co := item.(client.Correlation)
					return []string{co.UID, co.SourceUID, co.TargetUID, co.Label, co.Description}
				},
			})
		},
	}
}
