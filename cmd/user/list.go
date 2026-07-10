package user

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdUserList(f *cmdutil.Factory) *cobra.Command {
	var (
		query string
		page  int
		limit int
		all   bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long: `List users (requires server admin permissions). Use --all to fetch every page.

The output includes ID, Login, Email, Name, Admin status, and Disabled status.

Examples:
  # List all users
  grafana user list

  # Search users
  grafana user list -q "john"

  # Paginate
  grafana user list --page 2 --limit 50

  # Output as JSON
  grafana user list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			var users []client.User
			if all {
				for currentPage := 1; ; currentPage++ {
					result, err := c.ListUsers(cmd.Context(), query, client.PageParams{Page: currentPage, PerPage: limit})
					if err != nil {
						return err
					}
					users = append(users, result.Users...)
					if len(result.Users) == 0 ||
						(result.TotalCount > 0 && len(users) >= result.TotalCount) ||
						(limit > 0 && len(result.Users) < limit) {
						break
					}
				}
			} else {
				result, err := c.ListUsers(cmd.Context(), query, client.PageParams{Page: page, PerPage: limit})
				if err != nil {
					return err
				}
				users = result.Users
			}

			if len(users) == 0 && f.Resolved.Output == "table" {
				fmt.Fprintln(f.IOStreams.Out, "No users found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, users, &output.TableDef{
				Headers: []string{"ID", "Login", "Email", "Name", "Admin", "Disabled"},
				RowFunc: func(item interface{}) []string {
					u := item.(client.User)
					return []string{
						fmt.Sprintf("%d", u.ID),
						u.Login,
						u.Email,
						u.Name,
						fmt.Sprintf("%v", u.IsAdmin),
						fmt.Sprintf("%v", u.IsDisabled),
					}
				},
			})
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Search query")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)
	cmd.Flags().BoolVar(&all, "all", false, "Fetch all matching users across every page")

	return cmd
}
