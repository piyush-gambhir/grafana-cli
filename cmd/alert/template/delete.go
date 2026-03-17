package template

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdTemplateDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a notification template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete template %q?", name), confirm)
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

			if err := c.DeleteAlertTemplate(context.Background(), name); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Template %q deleted.\n", name)
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)

	return cmd
}
