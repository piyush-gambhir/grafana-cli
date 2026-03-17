package playlist

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdPlaylist returns the playlist parent command.
func NewCmdPlaylist(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "playlist",
		Short: "Manage playlists",
		Long: `Create, list, get, update, and delete dashboard playlists.

Playlists cycle through a sequence of dashboards at a configurable interval.`,
	}

	cmd.AddCommand(newCmdPlaylistList(f))
	cmd.AddCommand(newCmdPlaylistGet(f))
	cmd.AddCommand(newCmdPlaylistCreate(f))
	cmd.AddCommand(newCmdPlaylistUpdate(f))
	cmd.AddCommand(newCmdPlaylistDelete(f))

	return cmd
}
