package admin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAdminStats(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show Grafana server statistics",
		Long: `Display Grafana server usage statistics (requires admin permissions).

Shows counts of orgs, dashboards, datasources, users, alerts, etc.

Examples:
  # Show statistics
  grafana admin stats

  # Output as JSON
  grafana admin stats -o json`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAdminStats(context.Background())
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Organizations:       %d\n", result.Orgs)
				fmt.Fprintf(f.IOStreams.Out, "Dashboards:          %d\n", result.Dashboards)
				fmt.Fprintf(f.IOStreams.Out, "Datasources:         %d\n", result.Datasources)
				fmt.Fprintf(f.IOStreams.Out, "Users:               %d\n", result.Users)
				fmt.Fprintf(f.IOStreams.Out, "Active Users:        %d\n", result.ActiveUsers)
				fmt.Fprintf(f.IOStreams.Out, "Admins:              %d\n", result.Admins)
				fmt.Fprintf(f.IOStreams.Out, "Editors:             %d\n", result.Editors)
				fmt.Fprintf(f.IOStreams.Out, "Viewers:             %d\n", result.Viewers)
				fmt.Fprintf(f.IOStreams.Out, "Playlists:           %d\n", result.Playlists)
				fmt.Fprintf(f.IOStreams.Out, "Snapshots:           %d\n", result.Snapshots)
				fmt.Fprintf(f.IOStreams.Out, "Alerts:              %d\n", result.Alerts)
				fmt.Fprintf(f.IOStreams.Out, "Stars:               %d\n", result.Stars)
				fmt.Fprintf(f.IOStreams.Out, "Tags:                %d\n", result.Tags)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
