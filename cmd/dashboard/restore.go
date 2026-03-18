package dashboard

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdDashboardRestore(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <uid> <version>",
		Short: "Restore a dashboard to a specific version",
		Long: `Restore a dashboard to a specific historical version.

Use "grafana dashboard versions <uid>" to see available versions, then
restore to a previous version by specifying the dashboard UID and version
number.

Examples:
  # Restore dashboard to version 3
  grafana dashboard restore abc123 3

  # Check versions first, then restore
  grafana dashboard versions abc123
  grafana dashboard restore abc123 5`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			uid := args[0]
			version, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid version number: %s", args[1])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if err := c.RestoreDashboardVersion(context.Background(), uid, version); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Dashboard %q restored to version %d.\n", uid, version)
			return nil
		},
	}

	return cmd
}
