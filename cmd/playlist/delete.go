package playlist

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdPlaylistDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool
	var ifExists bool

	cmd := &cobra.Command{
		Use:         "delete <uid>",
		Short:       "Delete a playlist",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Permanently delete a playlist by its UID.

Examples:
  # Delete a playlist
  grafana playlist delete playlistUid

  # Delete without confirmation
  grafana playlist delete playlistUid --confirm

  # Delete idempotently (no error if not found)
  grafana playlist delete playlistUid --confirm --if-exists`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uid := args[0]

			ok, err := cmdutil.ConfirmAction(f.IOStreams.In, f.IOStreams.Out,
				fmt.Sprintf("Are you sure you want to delete playlist %q?", uid), confirm, f.NoInput)
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

			if err := c.DeletePlaylist(context.Background(), uid); err != nil {
				if ifExists && client.IsNotFound(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: playlist %q not found, skipping.\n", uid)
					}
					return nil
				}
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Playlist %q deleted.\n", uid)
			}
			return nil
		},
	}

	cmdutil.AddConfirmFlag(cmd, &confirm)
	cmdutil.AddIfExistsFlag(cmd, &ifExists)

	return cmd
}
