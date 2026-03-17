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

			// Get the dashboard to obtain its numeric ID.
			dash, err := c.GetDashboardByUID(context.Background(), uid)
			if err != nil {
				return err
			}

			dashID, ok := dash.Dashboard["id"].(float64)
			if !ok {
				return fmt.Errorf("could not determine dashboard ID")
			}

			if err := c.RestoreDashboardVersion(context.Background(), int64(dashID), version); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Dashboard %q restored to version %d.\n", uid, version)
			return nil
		},
	}

	return cmd
}
