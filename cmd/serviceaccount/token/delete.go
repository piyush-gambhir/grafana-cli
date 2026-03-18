package token

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdTokenDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <service-account-id> <token-id>",
		Short:       "Delete a service account token",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Delete an API token from a service account.

Examples:
  # Delete token 3 from service account 10
  grafana service-account token delete 10 3

  # Delete without confirmation
  grafana service-account token delete 10 3 --confirm

  # Delete idempotently (no error if not found)
  grafana service-account token delete 10 3 --confirm --if-exists`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			saID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			tokenID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token ID: %s", args[1])
			}

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete token %d?", tokenID), confirm, f.NoInput)
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

			if err := c.DeleteServiceAccountToken(context.Background(), saID, tokenID); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: token %d not found, skipping.\n", tokenID)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Token %d deleted.\n", tokenID)
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
