package snapshot

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdSnapshot returns the snapshot parent command.
func NewCmdSnapshot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage dashboard snapshots",
		Long: `Create, list, get, and delete dashboard snapshots.

Snapshots are point-in-time copies of dashboard data that can be shared
externally without requiring access to the Grafana instance.`,
	}

	cmd.AddCommand(newCmdSnapshotList(f))
	cmd.AddCommand(newCmdSnapshotGet(f))
	cmd.AddCommand(newCmdSnapshotCreate(f))
	cmd.AddCommand(newCmdSnapshotDelete(f))

	return cmd
}
