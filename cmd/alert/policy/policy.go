package policy

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdPolicy returns the notification policy parent command.
func NewCmdPolicy(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage notification policies",
		Long: `View, update, and reset the notification policy routing tree.

Notification policies define how alerts are routed to contact points,
including grouping, timing, and label matching.`,
	}

	cmd.AddCommand(newCmdPolicyGet(f))
	cmd.AddCommand(newCmdPolicyUpdate(f))
	cmd.AddCommand(newCmdPolicyReset(f))

	return cmd
}
