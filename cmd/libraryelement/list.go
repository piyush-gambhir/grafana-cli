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
		search string
		kind   int
		page   int
		limit  int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List library elements",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.ListLibraryElements(context.Background(), search, kind, client.PageParams{Page: page, PerPage: limit})
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

	cmd.Flags().StringVarP(&search, "search", "q", "", "Search string")
	cmd.Flags().IntVar(&kind, "kind", 0, "Kind (1=panel, 2=variable)")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
