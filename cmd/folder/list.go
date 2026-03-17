package folder

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdFolderList(f *cmdutil.Factory) *cobra.Command {
	var (
		page  int
		limit int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all folders",
		Long: `List all folders in the current organization.

The output includes ID, UID, and Title. Folders are used to organize
dashboards and alert rules. Use folder UIDs when filtering dashboards
or placing new dashboards.

Examples:
  # List all folders
  grafana folder list

  # Paginate results
  grafana folder list --page 1 --limit 50

  # Output as JSON
  grafana folder list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListFolders(context.Background(), client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No folders found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "UID", "Title"},
				RowFunc: func(item interface{}) []string {
					fl := item.(client.Folder)
					return []string{
						fmt.Sprintf("%d", fl.ID),
						fl.UID,
						fl.Title,
					}
				},
			})
		},
	}

	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
