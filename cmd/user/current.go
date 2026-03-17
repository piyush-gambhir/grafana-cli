package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdUserCurrent(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Short:   "Show the current authenticated user",
		Long: `Display information about the currently authenticated user.

Examples:
  # Show current user
  grafana user current

  # Alias
  grafana user whoami

  # Output as JSON
  grafana user current -o json`,
		Aliases: []string{"whoami"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetCurrentUser(context.Background())
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:    %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Login: %s\n", result.Login)
				fmt.Fprintf(f.IOStreams.Out, "Email: %s\n", result.Email)
				fmt.Fprintf(f.IOStreams.Out, "Name:  %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Admin: %v\n", result.IsAdmin)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
