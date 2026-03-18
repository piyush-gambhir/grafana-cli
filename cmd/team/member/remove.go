package member

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdMemberRemove(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:         "remove <team-id> <user-id>",
		Short:       "Remove a member from a team",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Remove a user from a team by team ID and user ID.

Examples:
  # Remove user 10 from team 5
  grafana team member remove 5 10

  # Remove without confirmation
  grafana team member remove 5 10 --confirm`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid team ID: %s", args[0])
			}

			userID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[1])
			}

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Remove user %d from team %d?", userID, teamID), confirm)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(f.IOStreams.Out, "Aborted.")
				return nil
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if err := c.RemoveTeamMember(context.Background(), teamID, userID); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "User %d removed from team %d.\n", userID, teamID)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
