package admin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdAdminReload(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload <resource>",
		Short: "Reload provisioned resources",
		Long: `Reload provisioned resources. Supported resources:
  dashboards, datasources, plugins, access-control, alerting`,
		ValidArgs: []string{"dashboards", "datasources", "plugins", "access-control", "alerting"},
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			resource := args[0]
			ctx := context.Background()

			switch resource {
			case "dashboards":
				err = c.ReloadDashboards(ctx)
			case "datasources":
				err = c.ReloadDatasources(ctx)
			case "plugins":
				err = c.ReloadPlugins(ctx)
			case "access-control":
				err = c.ReloadAccessControl(ctx)
			case "alerting":
				err = c.ReloadAlerting(ctx)
			default:
				return fmt.Errorf("unknown resource: %s (use dashboards, datasources, plugins, access-control, or alerting)", resource)
			}

			if err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Provisioned %s reloaded.\n", resource)
			return nil
		},
	}

	return cmd
}
