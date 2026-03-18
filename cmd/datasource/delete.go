package datasource

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdDatasourceDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <uid>",
		Short:       "Delete a datasource",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Permanently delete a datasource by its UID.

This action cannot be undone. Dashboards using this datasource will show
errors. You will be prompted for confirmation unless --confirm is set.

Examples:
  # Delete a datasource (with confirmation prompt)
  grafana datasource delete P1234

  # Delete without confirmation
  grafana datasource delete P1234 --confirm

  # Delete idempotently (no error if not found)
  grafana datasource delete P1234 --confirm --if-exists`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uid := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete datasource %q?", uid), confirm, f.NoInput)
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

			if err := c.DeleteDatasourceByUID(context.Background(), uid); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: datasource %q not found, skipping.\n", uid)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Datasource %q deleted.\n", uid)
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
