package correlation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdCorrelationDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <source-uid> <correlation-uid>",
		Short:       "Delete a correlation",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Delete a correlation by source datasource UID and correlation UID.

Examples:
  # Delete a correlation
  grafana correlation delete sourceUid corrUid

  # Delete without confirmation
  grafana correlation delete sourceUid corrUid --confirm

  # Delete idempotently (no error if not found)
  grafana correlation delete sourceUid corrUid --confirm --if-exists`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete correlation %q?", args[1]), confirm, f.NoInput)
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

			if err := c.DeleteCorrelation(context.Background(), args[0], args[1]); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: correlation %q not found, skipping.\n", args[1])
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Correlation %q deleted.\n", args[1])
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
