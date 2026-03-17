package user

import (
	"context"
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
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.ListUsers(context.Background(), query, client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(result.Users) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No users found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.Users, &output.TableDef{
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

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
