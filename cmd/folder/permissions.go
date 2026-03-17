package folder

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdFolderPermissions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage folder permissions",
	}

	cmd.AddCommand(newCmdFolderPermissionsGet(f))
	cmd.AddCommand(newCmdFolderPermissionsUpdate(f))

	return cmd
}

func newCmdFolderPermissionsGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get folder permissions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			perms, err := c.GetFolderPermissions(context.Background(), args[0])
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
					p := item.(client.FolderPermission)
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

func newCmdFolderPermissionsUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <uid>",
		Short: "Update folder permissions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.FolderPermissionsUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateFolderPermissions(context.Background(), args[0], req); err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Folder permissions updated.")
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
