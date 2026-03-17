package correlation

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdCorrelation returns the correlation parent command.
func NewCmdCorrelation(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "correlation",
		Short: "Manage correlations",
		Long: `Create, list, get, update, and delete datasource correlations.

Correlations link data between datasources, enabling navigation from
one datasource query result to another.`,
	}

	cmd.AddCommand(newCmdCorrelationList(f))
	cmd.AddCommand(newCmdCorrelationGet(f))
	cmd.AddCommand(newCmdCorrelationCreate(f))
	cmd.AddCommand(newCmdCorrelationUpdate(f))
	cmd.AddCommand(newCmdCorrelationDelete(f))

	return cmd
}
