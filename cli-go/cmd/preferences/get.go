package preferences

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/output"
)

func newCmdPreferencesGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current user preferences",
		Long: `Display the authenticated user's preferences.

Examples:
  # Get preferences
  grafana preferences get

  # Get as JSON
  grafana preferences get -o json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetPreferences(cmd.Context())
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Theme:     %s\n", result.Theme)
				fmt.Fprintf(f.IOStreams.Out, "Timezone:  %s\n", result.Timezone)
				fmt.Fprintf(f.IOStreams.Out, "WeekStart: %s\n", result.WeekStart)
				fmt.Fprintf(f.IOStreams.Out, "Language:  %s\n", result.Language)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
