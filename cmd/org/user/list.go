package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdOrgUserList(f *cmdutil.Factory) *cobra.Command {
	var (
		role  string
		query string
	)

	cmd := &cobra.Command{
		Use:     "list <org-id>",
		Short:   "List users in an organization",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Long: `List all users that belong to a specific organization.

The output includes User ID, Login, Email, Name, and Role. Results can be
filtered by role (Viewer, Editor, Admin) and searched by login, email, or
name using the --query flag.

The --role filter is case-insensitive and must be one of: Viewer, Editor,
or Admin. The --query filter performs a case-insensitive substring match
against the user's login, email, or name fields.

Examples:
  # List all users in org 1
  grafana org user list 1

  # Filter by role
  grafana org user list 1 --role Admin

  # Search by name or email
  grafana org user list 1 --query "john"

  # Combine role and query filters
  grafana org user list 1 --role Editor --query "dev"

  # Output as JSON
  grafana org user list 1 -o json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			orgID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid org ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListOrgUsers(context.Background(), orgID)
			if err != nil {
				return err
			}

			// Apply client-side filters.
			var filtered []client.OrgUser
			for _, u := range results {
				if role != "" && !strings.EqualFold(u.Role, role) {
					continue
				}
				if query != "" {
					q := strings.ToLower(query)
					if !strings.Contains(strings.ToLower(u.Login), q) &&
						!strings.Contains(strings.ToLower(u.Email), q) &&
						!strings.Contains(strings.ToLower(u.Name), q) {
						continue
					}
				}
				filtered = append(filtered, u)
			}

			if len(filtered) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No users found in this organization.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, filtered, &output.TableDef{
				Headers: []string{"User ID", "Login", "Email", "Name", "Role"},
				RowFunc: func(item interface{}) []string {
					u := item.(client.OrgUser)
					return []string{
						fmt.Sprintf("%d", u.UserID),
						u.Login,
						u.Email,
						u.Name,
						u.Role,
					}
				},
			})
		},
	}

	cmd.Flags().StringVar(&role, "role", "", "Filter by role (Viewer, Editor, Admin)")
	cmd.Flags().StringVar(&query, "query", "", "Search by login, email, or name")

	return cmd
}
