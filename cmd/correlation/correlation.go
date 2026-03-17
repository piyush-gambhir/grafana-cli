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
	}

	cmd.AddCommand(newCmdCorrelationList(f))
	cmd.AddCommand(newCmdCorrelationGet(f))
	cmd.AddCommand(newCmdCorrelationCreate(f))
	cmd.AddCommand(newCmdCorrelationUpdate(f))
	cmd.AddCommand(newCmdCorrelationDelete(f))

	return cmd
}
