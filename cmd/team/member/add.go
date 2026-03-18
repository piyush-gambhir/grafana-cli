package member

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdMemberAdd(f *cmdutil.Factory) *cobra.Command {
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "add <team-id> <user-id>",
		Short:       "Add a member to a team",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Add a user to a team by team ID and user ID.

Examples:
  # Add user 10 to team 5
  grafana team member add 5 10

  # Add idempotently (no error if already a member)
  grafana team member add 5 10 --if-not-exists`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid team ID: %s", args[0])
			}

			userID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[1])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			req := client.TeamMemberAddRequest{UserID: userID}
			if err := c.AddTeamMember(context.Background(), teamID, req); err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: user %d already in team %d, skipping.\n", userID, teamID)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "User %d added to team %d.\n", userID, teamID)
			}
			return nil
		},
	}

	cmdutil.AddIfNotExistsFlag(cmd, &ifNotExists)

	return cmd
}
