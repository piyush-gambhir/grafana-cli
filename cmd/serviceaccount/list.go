package serviceaccount

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdServiceAccountList(f *cmdutil.Factory) *cobra.Command {
	var (
		query string
		page  int
		limit int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List service accounts",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.ListServiceAccounts(context.Background(), query, client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(result.ServiceAccounts) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No service accounts found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.ServiceAccounts, &output.TableDef{
				Headers: []string{"ID", "Name", "Login", "Role", "Tokens", "Disabled"},
				RowFunc: func(item interface{}) []string {
					sa := item.(client.ServiceAccount)
					return []string{
						fmt.Sprintf("%d", sa.ID),
						sa.Name,
						sa.Login,
						sa.Role,
						fmt.Sprintf("%d", sa.Tokens),
						fmt.Sprintf("%v", sa.IsDisabled),
					}
				},
			})
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
