package rule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdRuleDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <uid>",
		Short: "Delete an alert rule",
		Long: `Permanently delete an alert rule by its UID.

Examples:
  # Delete an alert rule (with confirmation)
  grafana alert rule delete ruleUid123

  # Delete without confirmation
  grafana alert rule delete ruleUid123 --confirm`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uid := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete alert rule %q?", uid), confirm)
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

			if err := c.DeleteAlertRule(context.Background(), uid); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Alert rule %q deleted.\n", uid)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
