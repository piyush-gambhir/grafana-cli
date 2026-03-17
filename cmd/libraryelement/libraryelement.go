package libraryelement

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdLibraryElement returns the library element parent command.
func NewCmdLibraryElement(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "library-element",
		Short:   "Manage library elements",
		Long: `Create, list, get, update, delete library elements, and view connections.

Library elements are reusable panels and variables that can be shared
across multiple dashboards. When updated, changes propagate to all
connected dashboards.`,
		Aliases: []string{"le"},
	}

	cmd.AddCommand(newCmdLibraryElementList(f))
	cmd.AddCommand(newCmdLibraryElementGet(f))
	cmd.AddCommand(newCmdLibraryElementCreate(f))
	cmd.AddCommand(newCmdLibraryElementUpdate(f))
	cmd.AddCommand(newCmdLibraryElementDelete(f))
	cmd.AddCommand(newCmdLibraryElementConnections(f))

	return cmd
}
