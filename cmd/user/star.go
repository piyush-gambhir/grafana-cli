package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdUserStar(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "star",
		Short: "Star or unstar a dashboard",
		Long: `Star or unstar a dashboard for the current user.`,
	}

	cmd.AddCommand(newCmdUserStarAdd(f))
	cmd.AddCommand(newCmdUserStarRemove(f))

	return cmd
}

func newCmdUserStarAdd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:         "add <dashboard-id>",
		Short:       "Star a dashboard",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Star a dashboard by its numeric ID.

Examples:
  # Star dashboard 42
  grafana user star add 42`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid dashboard ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if err := c.StarDashboard(context.Background(), id); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Dashboard %d starred.\n", id)
			}
			return nil
		},
	}
}

func newCmdUserStarRemove(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:         "remove <dashboard-id>",
		Short:       "Unstar a dashboard",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Remove a star from a dashboard by its numeric ID.

Examples:
  # Unstar dashboard 42
  grafana user star remove 42`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid dashboard ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if err := c.UnstarDashboard(context.Background(), id); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Dashboard %d unstarred.\n", id)
			}
			return nil
		},
	}
}
