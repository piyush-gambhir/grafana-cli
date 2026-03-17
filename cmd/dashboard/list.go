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
		Long: `Search and list dashboards with optional filters.

The output includes UID, Title, Folder, Tags, and URL for each dashboard.
Results can be filtered by search query, tag, and folder UID.

The --query flag performs a full-text search on dashboard titles. The --tag
flag filters by a specific tag. The --folder flag filters by folder UID
(use "grafana folder list" to find folder UIDs).

Examples:
  # List all dashboards
  grafana dashboard list

  # Search by title
  grafana dashboard list -q "production"

  # Filter by tag
  grafana dashboard list --tag monitoring

  # Filter by folder UID
  grafana dashboard list --folder abc123

  # Paginate results
  grafana dashboard list --page 2 --limit 50

  # Output as JSON for scripting
  grafana dashboard list -o json

  # Output as YAML
  grafana dashboard list -o yaml`,
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
