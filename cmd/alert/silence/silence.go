package silence

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdSilence returns the silence parent command.
func NewCmdSilence(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "silence",
		Short: "Manage alert silences",
		Long: `Create, list, get, and delete alert silences.

Silences suppress alert notifications that match specific label matchers
during a defined time window.`,
	}

	cmd.AddCommand(newCmdSilenceList(f))
	cmd.AddCommand(newCmdSilenceGet(f))
	cmd.AddCommand(newCmdSilenceCreate(f))
	cmd.AddCommand(newCmdSilenceDelete(f))

	return cmd
}
