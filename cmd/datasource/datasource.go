package datasource

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdDatasource returns the datasource parent command.
func NewCmdDatasource(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasource",
		Short:   "Manage datasources",
		Long:    "Create, read, update, and delete Grafana datasources.",
		Aliases: []string{"ds"},
	}

	cmd.AddCommand(newCmdDatasourceList(f))
	cmd.AddCommand(newCmdDatasourceGet(f))
	cmd.AddCommand(newCmdDatasourceCreate(f))
	cmd.AddCommand(newCmdDatasourceUpdate(f))
	cmd.AddCommand(newCmdDatasourceDelete(f))

	return cmd
}
