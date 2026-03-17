package org

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdOrgCurrent(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current organization",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetCurrentOrg(context.Background())
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
