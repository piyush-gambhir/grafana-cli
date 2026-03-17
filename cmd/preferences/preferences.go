package preferences

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdPreferences returns the preferences parent command.
func NewCmdPreferences(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "preferences",
		Short:   "Manage user preferences",
		Long: `View and update the current user's preferences (theme, timezone, etc.).`,
		Aliases: []string{"prefs"},
	}

	cmd.AddCommand(newCmdPreferencesGet(f))
	cmd.AddCommand(newCmdPreferencesUpdate(f))

	return cmd
}
