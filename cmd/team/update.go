package team

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdTeamUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a team",
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

			var req client.TeamUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateTeam(context.Background(), id, req); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Team %d updated.\n", id)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
