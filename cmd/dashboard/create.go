package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		file        string
		folderUID   string
		overwrite   bool
		message     string
		ifNotExists bool
	)

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create a dashboard",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new dashboard from a JSON or YAML file.

The input file must contain a valid Grafana dashboard model. The "id"
field is automatically removed to ensure a new dashboard is created.
Use --folder to place the dashboard in a specific folder. Use --overwrite
to replace an existing dashboard with the same UID.

The --message flag sets a commit message for the dashboard version history.

Examples:
  # Create a dashboard from a JSON file
  grafana dashboard create -f dashboard.json

  # Create in a specific folder
  grafana dashboard create -f dashboard.json --folder folderUid123

  # Create from YAML
  grafana dashboard create -f dashboard.yaml

  # Overwrite if UID already exists
  grafana dashboard create -f dashboard.json --overwrite

  # Read from stdin
  cat dashboard.json | grafana dashboard create -f -

  # Set a version history message
  grafana dashboard create -f dashboard.json -m "Initial version"

  # Create idempotently (no error if already exists)
  grafana dashboard create -f dashboard.json --if-not-exists`,
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

			// Remove id to ensure creation.
			delete(dashboardData, "id")

			req := client.DashboardCreateRequest{
				Dashboard: dashboardData,
				FolderUID: folderUID,
				Overwrite: overwrite,
				Message:   message,
			}

			result, err := c.CreateDashboard(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: dashboard already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Dashboard created successfully.\n")
					fmt.Fprintf(f.IOStreams.Out, "  UID:     %s\n", result.UID)
					fmt.Fprintf(f.IOStreams.Out, "  URL:     %s\n", result.URL)
					fmt.Fprintf(f.IOStreams.Out, "  Version: %d\n", result.Version)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)
	cmd.Flags().StringVar(&folderUID, "folder", "", "Folder UID to place the dashboard in")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing dashboard with same UID")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Commit message for version history")
	cmdutil.AddIfNotExistsFlag(cmd, &ifNotExists)

	return cmd
}
