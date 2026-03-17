package libraryelement

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdLibraryElementConnections(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "connections <uid>",
		Short: "List dashboards connected to a library element",
		Long: `List all dashboards that use a specific library element.

Examples:
  # List connections
  grafana library-element connections leUid123

  # Output as JSON
  grafana library-element connections leUid123 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetLibraryElementConnections(context.Background(), args[0])
			if err != nil {
				return err
			}

			if len(result.Result) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No connections found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.Result, &output.TableDef{
				Headers: []string{"ID", "Element ID", "Connection ID", "Kind", "Created"},
				RowFunc: func(item interface{}) []string {
					conn := item.(client.LibraryElementConnection)
					return []string{
						fmt.Sprintf("%d", conn.ID),
						fmt.Sprintf("%d", conn.ElementID),
						fmt.Sprintf("%d", conn.ConnectionID),
						fmt.Sprintf("%d", conn.Kind),
						conn.Created,
					}
				},
			})
		},
	}
}
