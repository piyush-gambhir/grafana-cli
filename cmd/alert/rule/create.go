package rule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdRuleCreate(f *cmdutil.Factory) *cobra.Command {
	var file string
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create an alert rule",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new alert rule from a JSON or YAML file.

The file must contain a valid alert rule definition including title,
folderUID, ruleGroup, condition, data queries, and for duration.

Examples:
  # Create an alert rule
  grafana alert rule create -f rule.json

  # Create from YAML
  grafana alert rule create -f rule.yaml

  # Create idempotently (no error if already exists)
  grafana alert rule create -f rule.json --if-not-exists`,
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

			result, err := c.CreateAlertRule(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: alert rule already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Alert rule created: %s (UID: %s)\n", result.Title, result.UID)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)
	cmdutil.AddIfNotExistsFlag(cmd, &ifNotExists)

	return cmd
}
