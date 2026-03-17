package admin

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAdminSettings(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Show Grafana server settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAdminSettings(context.Background())
			if err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
