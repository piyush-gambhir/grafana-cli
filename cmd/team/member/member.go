package member

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdMember returns the team member parent command.
func NewCmdMember(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage team members",
		Long: `List, add, and remove members from a team.`,
	}

	cmd.AddCommand(newCmdMemberList(f))
	cmd.AddCommand(newCmdMemberAdd(f))
	cmd.AddCommand(newCmdMemberRemove(f))

	return cmd
}
