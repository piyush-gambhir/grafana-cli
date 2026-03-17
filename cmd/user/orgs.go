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

func newCmdUserOrgs(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "orgs <user-id>",
		Short: "List organizations a user belongs to",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.GetUserOrgs(context.Background(), userID)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "User is not in any organizations.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"Org ID", "Name", "Role"},
				RowFunc: func(item interface{}) []string {
					o := item.(client.UserOrg)
					return []string{fmt.Sprintf("%d", o.OrgID), o.Name, o.Role}
				},
			})
		},
	}
}
