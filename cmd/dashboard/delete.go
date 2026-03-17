package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdDashboardDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <uid>",
		Short: "Delete a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uid := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete dashboard %q?", uid), confirm)
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

			if err := c.DeleteDashboard(context.Background(), uid); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Dashboard %q deleted.\n", uid)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
