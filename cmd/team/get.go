package team

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTeamGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a team by ID",
		Long: `Retrieve a single team by its numeric ID.

Examples:
  # Get team details
  grafana team get 5

  # Get as JSON
  grafana team get 5 -o json`,
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

			result, err := c.GetTeam(context.Background(), id)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:      %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Name:    %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Email:   %s\n", result.Email)
				fmt.Fprintf(f.IOStreams.Out, "Members: %d\n", result.MemberCount)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
