package rule

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdRule returns the alert rule parent command.
func NewCmdRule(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage alert rules",
	}

	cmd.AddCommand(newCmdRuleList(f))
	cmd.AddCommand(newCmdRuleGet(f))
	cmd.AddCommand(newCmdRuleCreate(f))
	cmd.AddCommand(newCmdRuleUpdate(f))
	cmd.AddCommand(newCmdRuleDelete(f))

	return cmd
}
