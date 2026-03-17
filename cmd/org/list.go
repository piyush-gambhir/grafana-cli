package org

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdOrgList(f *cmdutil.Factory) *cobra.Command {
	var (
		page  int
		limit int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List organizations",
		Long: `List all organizations (requires server admin permissions).

The output includes ID and Name for each organization.

Examples:
  # List all organizations
  grafana org list

  # Paginate
  grafana org list --page 1 --limit 50

  # Output as JSON
  grafana org list -o json`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListOrgs(context.Background(), client.PageParams{Page: page, PerPage: limit})
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No organizations found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Name"},
				RowFunc: func(item interface{}) []string {
					o := item.(client.Org)
					return []string{fmt.Sprintf("%d", o.ID), o.Name}
				},
			})
		},
	}

	cmdutil.AddPaginationFlags(cmd, &page, &limit)

	return cmd
}
