package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdOrgUserAdd(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "add <org-id>",
		Short: "Add a user to an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.OrgUserAddRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.AddOrgUser(context.Background(), orgID, req); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "User %q added to organization %d with role %s.\n", req.LoginOrEmail, orgID, req.Role)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
