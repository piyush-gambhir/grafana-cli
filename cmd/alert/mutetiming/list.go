package mutetiming

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdMuteTimingList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List mute timings",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListMuteTimings(context.Background())
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No mute timings found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"Name", "Intervals", "Provenance"},
				RowFunc: func(item interface{}) []string {
					mt := item.(client.MuteTiming)
					return []string{mt.Name, fmt.Sprintf("%d", len(mt.TimeIntervals)), mt.Provenance}
				},
			})
		},
	}
}
