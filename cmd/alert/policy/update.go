package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdPolicyUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the notification policy tree",
		Long: `Replace the entire notification policy routing tree from a JSON or YAML file.

WARNING: This replaces the entire policy tree, not just parts of it.
Export the current tree first with "grafana alert policy get -o json".

Examples:
  # Update notification policy
  grafana alert policy update -f policy.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.NotificationPolicy
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateNotificationPolicy(context.Background(), req); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Notification policy updated.")
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
