package user

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdOrgUser returns the org user parent command.
func NewCmdOrgUser(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage organization users",
		Long: `List, add, update, and remove users in an organization.

Users can have one of three roles within an organization: Viewer, Editor,
or Admin.`,
	}

	cmd.AddCommand(newCmdOrgUserList(f))
	cmd.AddCommand(newCmdOrgUserAdd(f))
	cmd.AddCommand(newCmdOrgUserUpdate(f))
	cmd.AddCommand(newCmdOrgUserRemove(f))

	return cmd
}
