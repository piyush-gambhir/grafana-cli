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
	}

	cmd.AddCommand(newCmdPolicyGet(f))
	cmd.AddCommand(newCmdPolicyUpdate(f))
	cmd.AddCommand(newCmdPolicyReset(f))

	return cmd
}
