package org

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdOrgUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an organization",
		Long: `Update an organization's name from a JSON or YAML file.

Examples:
  # Update org 2
  grafana org update 2 -f org.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.OrgUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateOrg(context.Background(), id, req); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Organization %d updated.\n", id)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
