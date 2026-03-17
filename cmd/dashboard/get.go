package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a dashboard by UID",
		Long: `Retrieve a single dashboard by its UID.

In table mode, displays key fields: Title, UID, URL, Slug, Folder UID,
Version, and Starred status. In JSON/YAML mode, returns the full dashboard
model including all panels, templating, and metadata.

Use "grafana dashboard list" to find dashboard UIDs.

Examples:
  # Get dashboard details
  grafana dashboard get abc123

  # Get full dashboard JSON (for export/backup)
  grafana dashboard get abc123 -o json

  # Get dashboard as YAML
  grafana dashboard get abc123 -o yaml`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetDashboardByUID(context.Background(), args[0])
			if err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, &output.TableDef{
				Headers: []string{"Field", "Value"},
				RowFunc: func(item interface{}) []string {
					return nil
				},
			})
		},
	}

	// For table output of a single dashboard, we display key-value pairs.
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		c, err := f.Client()
		if err != nil {
			return err
		}

		result, err := c.GetDashboardByUID(context.Background(), args[0])
		if err != nil {
			return err
		}

		if f.Resolved.Output == "table" {
			title, _ := result.Dashboard["title"].(string)
			uid, _ := result.Dashboard["uid"].(string)
			fmt.Fprintf(f.IOStreams.Out, "Title:      %s\n", title)
			fmt.Fprintf(f.IOStreams.Out, "UID:        %s\n", uid)
			fmt.Fprintf(f.IOStreams.Out, "URL:        %s\n", result.Meta.URL)
			fmt.Fprintf(f.IOStreams.Out, "Slug:       %s\n", result.Meta.Slug)
			fmt.Fprintf(f.IOStreams.Out, "Folder UID: %s\n", result.Meta.FolderUID)
			fmt.Fprintf(f.IOStreams.Out, "Version:    %d\n", result.Meta.Version)
			fmt.Fprintf(f.IOStreams.Out, "Starred:    %v\n", result.Meta.IsStarred)
			return nil
		}

		return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
	}

	return cmd
}
