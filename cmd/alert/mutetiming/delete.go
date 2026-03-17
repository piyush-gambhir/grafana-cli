package mutetiming

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdMuteTimingDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a mute timing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete mute timing %q?", name), confirm)
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

			if err := c.DeleteMuteTiming(context.Background(), name); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Mute timing %q deleted.\n", name)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
