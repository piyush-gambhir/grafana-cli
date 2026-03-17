package org

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdOrgGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an organization by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetOrg(context.Background(), id)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:   %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Name: %s\n", result.Name)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
