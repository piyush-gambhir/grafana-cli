package dashboard

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDashboardPermissions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage dashboard permissions",
		Long: `View and update access control permissions for a specific dashboard.

Subcommands:
  get    - Retrieve current permissions
  update - Replace permissions from a JSON/YAML file`,
	}

	cmd.AddCommand(newCmdDashboardPermissionsGet(f))
	cmd.AddCommand(newCmdDashboardPermissionsUpdate(f))

	return cmd
}

func newCmdDashboardPermissionsGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get dashboard permissions",
		Long: `Retrieve the current permissions for a dashboard.

The output includes User, Team, Role, and Permission level for each
access control entry.

Examples:
  # Get permissions for a dashboard
  grafana dashboard permissions get abc123

  # Output as JSON
  grafana dashboard permissions get abc123 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			perms, err := c.GetDashboardPermissions(context.Background(), args[0])
			if err != nil {
				return err
			}

			if len(perms) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No permissions found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, perms, &output.TableDef{
				Headers: []string{"User", "Team", "Role", "Permission"},
				RowFunc: func(item interface{}) []string {
					p := item.(client.DashboardPermission)
					user := p.UserLogin
					if user == "" {
						user = "-"
					}
					teamName := p.Team
					if teamName == "" {
						teamName = "-"
					}
					role := p.Role
					if role == "" {
						role = "-"
					}
					return []string{user, teamName, role, p.PermissionName}
				},
			})
		},
	}
}

func newCmdDashboardPermissionsUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <uid>",
		Short:       "Update dashboard permissions",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update dashboard permissions from a JSON or YAML file.

The file should contain a permissions object with an "items" array. Each
item specifies a user, team, or role and their permission level.

Examples:
  # Update permissions from file
  grafana dashboard permissions update abc123 -f perms.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.DashboardPermissionsUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateDashboardPermissions(context.Background(), args[0], req); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Dashboard permissions updated.")
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
