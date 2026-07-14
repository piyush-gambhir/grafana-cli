package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/output"
)

func newCmdTemplateList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List notification templates",
		Long: `List all notification templates.

The output includes Name and Provenance.

Examples:
  # List all templates
  grafana alert template list

  # Output as JSON
  grafana alert template list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListAlertTemplates(cmd.Context())
			if err != nil {
				return err
			}

			if len(results) == 0 && f.Resolved.Output == "table" {
				fmt.Fprintln(f.IOStreams.Out, "No notification templates found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"Name", "Provenance"},
				RowFunc: func(item interface{}) []string {
					t := item.(client.AlertTemplate)
					return []string{t.Name, t.Provenance}
				},
			})
		},
	}
}
