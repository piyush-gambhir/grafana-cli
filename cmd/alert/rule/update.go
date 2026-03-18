package rule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdRuleUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <uid>",
		Short:       "Update an alert rule",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update an existing alert rule from a JSON or YAML file.

Examples:
  # Update an alert rule
  grafana alert rule update ruleUid123 -f updated-rule.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.AlertRule
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.UpdateAlertRule(context.Background(), args[0], req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Alert rule updated: %s (UID: %s)\n", result.Title, result.UID)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
