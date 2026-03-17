package member

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdMemberList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <team-id>",
		Short:   "List team members",
		Long: `List all members of a team.

Examples:
  # List members of team 5
  grafana team member list 5

  # Output as JSON
  grafana team member list 5 -o json`,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid team ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListTeamMembers(context.Background(), teamID)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No members found in this team.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"User ID", "Login", "Email", "Name"},
				RowFunc: func(item interface{}) []string {
					m := item.(client.TeamMember)
					return []string{
						fmt.Sprintf("%d", m.UserID),
						m.Login,
						m.Email,
						m.Name,
					}
				},
			})
		},
	}
}
