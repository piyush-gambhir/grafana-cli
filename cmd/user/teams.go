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

func newCmdUserTeams(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "teams <user-id>",
		Short: "List teams a user belongs to",
		Long: `List all teams a user is a member of (requires server admin).

Examples:
  # List teams for user 5
  grafana user teams 5`,
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

			results, err := c.GetUserTeams(context.Background(), userID)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "User is not in any teams.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Name", "Email", "Members"},
				RowFunc: func(item interface{}) []string {
					t := item.(client.UserTeam)
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
}
