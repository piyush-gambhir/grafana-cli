package org

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdOrgSwitch(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "switch <org-id>",
		Short: "Switch the active organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if err := c.SwitchOrg(context.Background(), orgID); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Switched to organization %d.\n", orgID)
			return nil
		},
	}
}
