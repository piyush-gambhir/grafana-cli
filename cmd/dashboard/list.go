package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardList(f *cmdutil.Factory) *cobra.Command {
	var (
		query     string
		tag       string
		folderUID string
		page      int
		limit     int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		Long:  "Search and list dashboards with optional filters.",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.SearchDashboards(context.Background(), query, tag, folderUID, client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No dashboards found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"UID", "Title", "Folder", "Tags", "URL"},
				RowFunc: func(item interface{}) []string {
					d := item.(client.DashboardSearchResult)
					return []string{
						d.UID,
						d.Title,
						d.FolderTitle,
						strings.Join(d.Tags, ", "),
						d.URL,
					}
				},
			})
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Filter by tag")
	cmd.Flags().StringVar(&folderUID, "folder", "", "Filter by folder UID")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
