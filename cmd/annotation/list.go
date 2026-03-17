package annotation

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAnnotationList(f *cmdutil.Factory) *cobra.Command {
	var (
		dashboardID int64
		panelID     int64
		from        int64
		to          int64
		tags        []string
		limit       int64
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List annotations",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListAnnotations(context.Background(), dashboardID, panelID, from, to, tags, limit)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No annotations found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Dashboard ID", "Text", "Tags", "Time"},
				RowFunc: func(item interface{}) []string {
					a := item.(client.Annotation)
					return []string{
						fmt.Sprintf("%d", a.ID),
						fmt.Sprintf("%d", a.DashboardID),
						a.Text,
						strings.Join(a.Tags, ", "),
						fmt.Sprintf("%d", a.Time),
					}
				},
			})
		},
	}

	cmd.Flags().Int64Var(&dashboardID, "dashboard-id", 0, "Filter by dashboard ID")
	cmd.Flags().Int64Var(&panelID, "panel-id", 0, "Filter by panel ID")
	cmd.Flags().Int64Var(&from, "from", 0, "Start time (epoch ms)")
	cmd.Flags().Int64Var(&to, "to", 0, "End time (epoch ms)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Filter by tags")
	cmd.Flags().Int64Var(&limit, "limit", 100, "Maximum number of annotations to return")

	return cmd
}
