package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdOrgUserList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <org-id>",
		Short:   "List users in an organization",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListOrgUsers(context.Background(), orgID)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No users found in this organization.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"User ID", "Login", "Email", "Name", "Role"},
				RowFunc: func(item interface{}) []string {
					u := item.(client.OrgUser)
					return []string{
						fmt.Sprintf("%d", u.UserID),
						u.Login,
						u.Email,
						u.Name,
						u.Role,
					}
				},
			})
		},
	}
}
