package datasource

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDatasourceList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all datasources",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListDatasources(context.Background())
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No datasources found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "UID", "Name", "Type", "URL", "Default"},
				RowFunc: func(item interface{}) []string {
					d := item.(client.Datasource)
					isDefault := "no"
					if d.IsDefault {
						isDefault = "yes"
					}
					return []string{
						fmt.Sprintf("%d", d.ID),
						d.UID,
						d.Name,
						d.Type,
						d.URL,
						isDefault,
					}
				},
			})
		},
	}

	return cmd
}
