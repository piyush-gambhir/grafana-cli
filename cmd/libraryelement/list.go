package libraryelement

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdLibraryElementList(f *cmdutil.Factory) *cobra.Command {
	var (
		search    string
		kind      int
		folderFilter string
		page      int
		limit     int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List library elements",
		Aliases: []string{"ls"},
		Long: `List library elements (reusable panels and variables) in the current organization.

The output includes UID, Name, Kind (panel or variable), Type, Folder,
and the number of connected dashboards.

The --kind flag accepts 1 for panels or 2 for variables. The --folder flag
filters by folder filter (accepts folder name for the Grafana API). The
--search flag performs a search on the element name.

Examples:
  # List all library elements
  grafana library-element list

  # List only panels
  grafana library-element list --kind 1

  # List only variables
  grafana library-element list --kind 2

  # Search by name
  grafana library-element list --search "CPU"

  # Filter by folder
  grafana library-element list --folder "General"

  # Paginate results
  grafana library-element list --page 2 --limit 20

  # Output as JSON
  grafana library-element list -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.ListLibraryElements(context.Background(), search, kind, folderFilter, client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(result.Result.Elements) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No library elements found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.Result.Elements, &output.TableDef{
				Headers: []string{"UID", "Name", "Kind", "Type", "Folder", "Connections"},
				RowFunc: func(item interface{}) []string {
					le := item.(client.LibraryElement)
					kindStr := "panel"
					if le.Kind == 2 {
						kindStr = "variable"
					}
					return []string{
						le.UID,
						le.Name,
						kindStr,
						le.Type,
						le.Meta.FolderName,
						fmt.Sprintf("%d", le.Meta.ConnectedDashboards),
					}
				},
			})
		},
	}

	cmd.Flags().StringVar(&search, "search", "", "Search string for element name")
	cmd.Flags().IntVar(&kind, "kind", 0, "Kind filter (1=panel, 2=variable)")
	cmd.Flags().StringVar(&folderFilter, "folder", "", "Filter by folder name")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
