package folder

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdFolder returns the folder parent command.
func NewCmdFolder(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "folder",
		Short: "Manage folders",
		Long:  "Create, read, update, and delete Grafana folders.",
	}

	cmd.AddCommand(newCmdFolderList(f))
	cmd.AddCommand(newCmdFolderGet(f))
	cmd.AddCommand(newCmdFolderCreate(f))
	cmd.AddCommand(newCmdFolderUpdate(f))
	cmd.AddCommand(newCmdFolderDelete(f))
	cmd.AddCommand(newCmdFolderPermissions(f))

	return cmd
}
