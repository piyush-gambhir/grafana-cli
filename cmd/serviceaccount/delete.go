package serviceaccount

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdServiceAccountDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <id>",
		Short:       "Delete a service account",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Permanently delete a service account and all its tokens.

Examples:
  # Delete service account 10
  grafana service-account delete 10

  # Delete without confirmation
  grafana service-account delete 10 --confirm

  # Delete idempotently (no error if not found)
  grafana service-account delete 10 --confirm --if-exists`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete service account %d?", id), confirm, f.NoInput)
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

			if err := c.DeleteServiceAccount(context.Background(), id); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: service account %d not found, skipping.\n", id)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Service account %d deleted.\n", id)
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
