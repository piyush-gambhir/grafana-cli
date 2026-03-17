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
		Long: `Create, list, get, update, and delete Grafana Unified Alerting rules.

Alert rules define conditions that trigger alerts and notifications. Each
rule belongs to a rule group within a folder.`,
	}

	cmd.AddCommand(newCmdRuleList(f))
	cmd.AddCommand(newCmdRuleGet(f))
	cmd.AddCommand(newCmdRuleCreate(f))
	cmd.AddCommand(newCmdRuleUpdate(f))
	cmd.AddCommand(newCmdRuleDelete(f))

	return cmd
}
