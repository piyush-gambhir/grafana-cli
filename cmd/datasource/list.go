package datasource

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDatasourceList(f *cmdutil.Factory) *cobra.Command {
	var (
		dsType string
		name   string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all datasources",
		Aliases: []string{"ls"},
		Long: `List all configured datasources in the current organization.

The output includes ID, UID, Name, Type, URL, and whether the datasource
is the default. Results can be filtered by type or name.

The --type flag filters by the datasource plugin type identifier (e.g.
prometheus, elasticsearch, loki, graphite, influxdb, mysql, postgres,
cloudwatch, stackdriver, tempo, jaeger, zipkin, opentsdb, grafana-azure-monitor-datasource).

The --name flag performs a case-insensitive substring match on the datasource name.

Examples:
  # List all datasources
  grafana datasource list

  # List only Prometheus datasources
  grafana datasource list --type prometheus

  # Search datasources by name
  grafana datasource list --name "prod"

  # List datasources as JSON
  grafana datasource list -o json

  # Combine filters
  grafana datasource list --type loki --name "staging"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListDatasources(context.Background())
			if err != nil {
				return err
			}

			// Apply client-side filters since the Grafana API returns all datasources.
			var filtered []client.Datasource
			for _, d := range results {
				if dsType != "" && !strings.EqualFold(d.Type, dsType) {
					continue
				}
				if name != "" && !strings.Contains(strings.ToLower(d.Name), strings.ToLower(name)) {
					continue
				}
				filtered = append(filtered, d)
			}

			if len(filtered) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No datasources found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, filtered, &output.TableDef{
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

	cmd.Flags().StringVar(&dsType, "type", "", "Filter by datasource type (e.g. prometheus, elasticsearch, loki)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Filter by datasource name (case-insensitive substring match)")

	return cmd
}
