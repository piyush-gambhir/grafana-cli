package admin

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdAdmin returns the admin parent command.
func NewCmdAdmin(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Server administration commands",
		Long:  "Perform Grafana server administration tasks (requires admin permissions).",
	}

	cmd.AddCommand(newCmdAdminSettings(f))
	cmd.AddCommand(newCmdAdminStats(f))
	cmd.AddCommand(newCmdAdminReload(f))

	return cmd
}
