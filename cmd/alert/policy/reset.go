package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdPolicyReset(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:         "reset",
		Short:       "Reset the notification policy to defaults",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Reset the notification policy routing tree to Grafana defaults.

This removes all custom routes and restores the default receiver
configuration. You will be prompted for confirmation.

Examples:
  # Reset notification policy (with confirmation)
  grafana alert policy reset

  # Reset without confirmation
  grafana alert policy reset --confirm`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				"Are you sure you want to reset the notification policy to defaults?", confirm, f.NoInput)
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

			if err := c.ResetNotificationPolicy(context.Background()); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintln(f.IOStreams.Out, "Notification policy reset to defaults.")
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
