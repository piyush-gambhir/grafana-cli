package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdUserUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <id>",
		Short:       "Update a user",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update a user's profile from a JSON or YAML file.

Examples:
  # Update user 5
  grafana user update 5 -f user.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.UserUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateUser(context.Background(), id, req); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "User %d updated.\n", id)
			}
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
