package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdOrgUserRemove(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:         "remove <org-id> <user-id>",
		Short:       "Remove a user from an organization",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Remove a user from an organization.

Examples:
  # Remove user 5 from org 1
  grafana org user remove 1 5

  # Remove without confirmation
  grafana org user remove 1 5 --confirm`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			userID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[1])
			}

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Remove user %d from organization %d?", userID, orgID), confirm)
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

			if err := c.RemoveOrgUser(context.Background(), orgID, userID); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "User %d removed from organization %d.\n", userID, orgID)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
