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
	}

	cmd.AddCommand(newCmdSnapshotList(f))
	cmd.AddCommand(newCmdSnapshotGet(f))
	cmd.AddCommand(newCmdSnapshotCreate(f))
	cmd.AddCommand(newCmdSnapshotDelete(f))

	return cmd
}
