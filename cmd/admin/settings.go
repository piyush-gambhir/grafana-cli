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
		Long: `Display all Grafana server settings (requires admin permissions).

Returns the server configuration organized by section. Use JSON output
for machine-readable parsing.

Examples:
  # Show settings
  grafana admin settings

  # Output as JSON
  grafana admin settings -o json`,
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
