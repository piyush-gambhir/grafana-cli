package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		file      string
		folderUID string
		overwrite bool
		message   string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a dashboard",
		Long:  "Update an existing dashboard from a JSON or YAML file.",
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
				fmt.Fprintf(f.IOStreams.Out, "Dashboard updated successfully.\n")
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
	cmd.Flags().BoolVar(&overwrite, "overwrite", true, "Overwrite existing dashboard")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message for version history")

	return cmd
}
