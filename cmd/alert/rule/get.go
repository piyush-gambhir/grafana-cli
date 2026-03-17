package rule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdRuleGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get an alert rule by UID",
		Long: `Retrieve a single alert rule by its UID.

In table mode, displays UID, Title, Folder UID, Rule Group, Condition,
For duration, No Data state, and Exec Error state. In JSON/YAML mode,
returns the full rule definition including data queries.

Examples:
  # Get alert rule details
  grafana alert rule get ruleUid123

  # Get full JSON (for creating update payloads)
  grafana alert rule get ruleUid123 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			rule, err := c.GetAlertRule(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "UID:          %s\n", rule.UID)
				fmt.Fprintf(f.IOStreams.Out, "Title:        %s\n", rule.Title)
				fmt.Fprintf(f.IOStreams.Out, "Folder UID:   %s\n", rule.FolderUID)
				fmt.Fprintf(f.IOStreams.Out, "Rule Group:   %s\n", rule.RuleGroup)
				fmt.Fprintf(f.IOStreams.Out, "Condition:    %s\n", rule.Condition)
				fmt.Fprintf(f.IOStreams.Out, "For:          %s\n", rule.For)
				fmt.Fprintf(f.IOStreams.Out, "No Data:      %s\n", rule.NoDataState)
				fmt.Fprintf(f.IOStreams.Out, "Exec Error:   %s\n", rule.ExecErrState)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, rule, nil)
		},
	}
}
