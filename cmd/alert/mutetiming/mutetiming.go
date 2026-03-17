package mutetiming

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdMuteTiming returns the mute timing parent command.
func NewCmdMuteTiming(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mute-timing",
		Short:   "Manage mute timings",
		Long: `Create, list, get, update, and delete mute timings.

Mute timings define time periods during which alert notifications are
suppressed. They are referenced by name in notification policies.`,
		Aliases: []string{"mt"},
	}

	cmd.AddCommand(newCmdMuteTimingList(f))
	cmd.AddCommand(newCmdMuteTimingGet(f))
	cmd.AddCommand(newCmdMuteTimingCreate(f))
	cmd.AddCommand(newCmdMuteTimingUpdate(f))
	cmd.AddCommand(newCmdMuteTimingDelete(f))

	return cmd
}
