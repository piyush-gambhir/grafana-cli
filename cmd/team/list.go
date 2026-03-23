package team

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTeamList(f *cmdutil.Factory) *cobra.Command {
	var (
		query string
		page  int
		limit int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List teams",
		Long: `Search and list teams in the current organization.

The output includes ID, Name, Email, and Member count.

Examples:
  # List all teams
  grafana team list

  # Search teams by name
  grafana team list -q "backend"

  # Paginate results
  grafana team list --page 1 --limit 20

  # Output as JSON
  grafana team list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.ListTeams(context.Background(), query, client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(result.Teams) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No teams found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.Teams, &output.TableDef{
				Headers: []string{"ID", "Name", "Email", "Members"},
				RowFunc: func(item interface{}) []string {
					t := item.(client.Team)
					return []string{
						fmt.Sprintf("%d", t.ID),
						t.Name,
						t.Email,
						fmt.Sprintf("%d", t.MemberCount),
					}
				},
			})
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Search query")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
