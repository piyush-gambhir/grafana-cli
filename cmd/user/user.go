package user

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdUser returns the user parent command.
func NewCmdUser(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(newCmdUserList(f))
	cmd.AddCommand(newCmdUserGet(f))
	cmd.AddCommand(newCmdUserLookup(f))
	cmd.AddCommand(newCmdUserUpdate(f))
	cmd.AddCommand(newCmdUserOrgs(f))
	cmd.AddCommand(newCmdUserTeams(f))
	cmd.AddCommand(newCmdUserCurrent(f))
	cmd.AddCommand(newCmdUserStar(f))

	return cmd
}
