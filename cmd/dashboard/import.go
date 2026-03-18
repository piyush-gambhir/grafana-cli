package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardImport(f *cmdutil.Factory) *cobra.Command {
	var (
		file      string
		folderUID string
		overwrite bool
		message   string
	)

	cmd := &cobra.Command{
		Use:         "import",
		Short:       "Import a dashboard from a file",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Import a dashboard from a JSON or YAML file.

This is functionally equivalent to "dashboard create". The "id" field is
automatically stripped to allow clean imports. Use --folder to specify the
target folder and --overwrite to replace an existing dashboard.

Examples:
  # Import a dashboard
  grafana dashboard import -f exported-dashboard.json

  # Import into a specific folder
  grafana dashboard import -f dashboard.json --folder folderUid

  # Import and overwrite if exists
  grafana dashboard import -f dashboard.json --overwrite

  # Import from stdin
  cat dashboard.json | grafana dashboard import -f -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var dashboardData map[string]interface{}
			if err := cmdutil.UnmarshalInput(file, &dashboardData); err != nil {
				return err
			}

			// Remove id to allow import.
			delete(dashboardData, "id")

			req := client.DashboardCreateRequest{
				Dashboard: dashboardData,
				FolderUID: folderUID,
				Overwrite: overwrite,
				Message:   message,
			}

			result, err := c.CreateDashboard(context.Background(), req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Dashboard imported successfully.\n")
				fmt.Fprintf(f.IOStreams.Out, "  UID:     %s\n", result.UID)
				fmt.Fprintf(f.IOStreams.Out, "  URL:     %s\n", result.URL)
				fmt.Fprintf(f.IOStreams.Out, "  Version: %d\n", result.Version)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)
	cmd.Flags().StringVar(&folderUID, "folder", "", "Folder UID to place the dashboard in")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing dashboard with same UID")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message for version history")

	return cmd
}
