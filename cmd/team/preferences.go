package team

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTeamPreferences(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preferences",
		Short: "Manage team preferences",
	}

	cmd.AddCommand(newCmdTeamPreferencesGet(f))
	cmd.AddCommand(newCmdTeamPreferencesUpdate(f))

	return cmd
}

func newCmdTeamPreferencesGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <team-id>",
		Short: "Get team preferences",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid team ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			prefs, err := c.GetTeamPreferences(context.Background(), id)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Theme:     %s\n", prefs.Theme)
				fmt.Fprintf(f.IOStreams.Out, "Timezone:  %s\n", prefs.Timezone)
				fmt.Fprintf(f.IOStreams.Out, "WeekStart: %s\n", prefs.WeekStart)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, prefs, nil)
		},
	}
}

func newCmdTeamPreferencesUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <team-id>",
		Short: "Update team preferences",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid team ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var prefs client.TeamPreferences
			if err := cmdutil.UnmarshalInput(file, &prefs); err != nil {
				return err
			}

			if err := c.UpdateTeamPreferences(context.Background(), id, prefs); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Team %d preferences updated.\n", id)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
