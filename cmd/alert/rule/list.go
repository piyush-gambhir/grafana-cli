package rule

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdRuleList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List alert rules",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			rules, err := c.ListAlertRules(context.Background())
			if err != nil {
				return err
			}

			if len(rules) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No alert rules found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, rules, &output.TableDef{
				Headers: []string{"UID", "Title", "Folder UID", "Rule Group", "Condition", "For"},
				RowFunc: func(item interface{}) []string {
					r := item.(client.AlertRule)
					return []string{r.UID, r.Title, r.FolderUID, r.RuleGroup, r.Condition, r.For}
				},
			})
		},
	}
}
