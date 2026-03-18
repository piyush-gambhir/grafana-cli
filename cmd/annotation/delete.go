package annotation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdAnnotationDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <id>",
		Short:       "Delete an annotation",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Permanently delete an annotation by its numeric ID.

Examples:
  # Delete annotation 42
  grafana annotation delete 42

  # Delete without confirmation
  grafana annotation delete 42 --confirm

  # Delete idempotently (no error if not found)
  grafana annotation delete 42 --confirm --if-exists`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid annotation ID: %s", args[0])
			}

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete annotation %d?", id), confirm, f.NoInput)
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

			if err := c.DeleteAnnotation(context.Background(), id); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: annotation %d not found, skipping.\n", id)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Annotation %d deleted.\n", id)
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
