package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdOrgUserUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <org-id> <user-id>",
		Short: "Update a user's role in an organization",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			userID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[1])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.OrgUserUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateOrgUser(context.Background(), orgID, userID, req); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "User %d updated in organization %d.\n", userID, orgID)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
