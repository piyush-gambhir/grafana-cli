package org

import (
	"github.com/spf13/cobra"

	orguser "github.com/piyush-gambhir/grafana-cli/cmd/org/user"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdOrg returns the org parent command.
func NewCmdOrg(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Manage organizations",
	}

	cmd.AddCommand(newCmdOrgList(f))
	cmd.AddCommand(newCmdOrgGet(f))
	cmd.AddCommand(newCmdOrgCreate(f))
	cmd.AddCommand(newCmdOrgUpdate(f))
	cmd.AddCommand(newCmdOrgDelete(f))
	cmd.AddCommand(newCmdOrgCurrent(f))
	cmd.AddCommand(newCmdOrgSwitch(f))
	cmd.AddCommand(orguser.NewCmdOrgUser(f))

	return cmd
}
