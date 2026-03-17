package dashboard

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdDashboard returns the dashboard parent command.
func NewCmdDashboard(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Short:   "Manage dashboards",
		Long:    "Create, read, update, and delete Grafana dashboards.",
		Aliases: []string{"dash", "db"},
	}

	cmd.AddCommand(newCmdDashboardList(f))
	cmd.AddCommand(newCmdDashboardGet(f))
	cmd.AddCommand(newCmdDashboardCreate(f))
	cmd.AddCommand(newCmdDashboardUpdate(f))
	cmd.AddCommand(newCmdDashboardDelete(f))
	cmd.AddCommand(newCmdDashboardExport(f))
	cmd.AddCommand(newCmdDashboardImport(f))
	cmd.AddCommand(newCmdDashboardVersions(f))
	cmd.AddCommand(newCmdDashboardRestore(f))
	cmd.AddCommand(newCmdDashboardPermissions(f))

	return cmd
}
