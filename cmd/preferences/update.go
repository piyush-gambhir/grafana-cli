package preferences

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdPreferencesUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update",
		Short:       "Update current user preferences",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update the authenticated user's preferences from a JSON or YAML file.

Examples:
  # Update preferences
  grafana preferences update -f prefs.json

  # Example JSON: {"theme":"dark","timezone":"utc","weekStart":"monday"}`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.PreferencesUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdatePreferences(context.Background(), req); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintln(f.IOStreams.Out, "Preferences updated.")
			}
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
