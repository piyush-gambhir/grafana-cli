package rule

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdRuleList(f *cmdutil.Factory) *cobra.Command {
	var (
		folder string
		group  string
		limit  int
		page   int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List alert rules",
		Aliases: []string{"ls"},
		Long: `List all alert rules in the current organization.

The output includes UID, Title, Folder UID, Rule Group, Condition, and the
For duration. Results can be filtered by folder UID and rule group.

The Grafana provisioning API returns all rules. The --folder and --group
flags perform client-side filtering on the folder UID and rule group name
respectively. The --limit and --page flags control pagination of the
filtered results.

Examples:
  # List all alert rules
  grafana alert rule list

  # Filter by folder UID
  grafana alert rule list --folder abc123

  # Filter by rule group
  grafana alert rule list --group "High CPU"

  # Combine folder and group filters
  grafana alert rule list --folder abc123 --group "High CPU"

  # Paginate results
  grafana alert rule list --limit 10 --page 2

  # Output as JSON for scripting
  grafana alert rule list -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			rules, err := c.ListAlertRules(context.Background())
			if err != nil {
				return err
			}

			// Apply client-side filters.
			var filtered []client.AlertRule
			for _, r := range rules {
				if folder != "" && !strings.EqualFold(r.FolderUID, folder) {
					continue
				}
				if group != "" && r.RuleGroup != group {
					continue
				}
				filtered = append(filtered, r)
			}

			// Apply pagination.
			total := len(filtered)
			if limit > 0 {
				start := (page - 1) * limit
				if start >= total {
					filtered = nil
				} else {
					end := start + limit
					if end > total {
						end = total
					}
					filtered = filtered[start:end]
				}
			}

			if len(filtered) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No alert rules found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, filtered, &output.TableDef{
				Headers: []string{"UID", "Title", "Folder UID", "Rule Group", "Condition", "For"},
				RowFunc: func(item interface{}) []string {
					r := item.(client.AlertRule)
					return []string{r.UID, r.Title, r.FolderUID, r.RuleGroup, r.Condition, r.For}
				},
			})
		},
	}

	cmd.Flags().StringVar(&folder, "folder", "", "Filter by folder UID")
	cmd.Flags().StringVar(&group, "group", "", "Filter by rule group name")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of rules to return (0 = all)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number (used with --limit)")

	return cmd
}
