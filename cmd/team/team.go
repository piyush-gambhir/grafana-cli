package team

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cmd/team/member"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdTeam returns the team parent command.
func NewCmdTeam(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Manage teams",
	}

	cmd.AddCommand(newCmdTeamList(f))
	cmd.AddCommand(newCmdTeamGet(f))
	cmd.AddCommand(newCmdTeamCreate(f))
	cmd.AddCommand(newCmdTeamUpdate(f))
	cmd.AddCommand(newCmdTeamDelete(f))
	cmd.AddCommand(newCmdTeamPreferences(f))
	cmd.AddCommand(member.NewCmdMember(f))

	return cmd
}
