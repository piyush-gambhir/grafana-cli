package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdUserLookup(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "lookup <login-or-email>",
		Short: "Look up a user by login or email",
		Long: `Find a user by their login username or email address.

Examples:
  # Look up by login
  grafana user lookup admin

  # Look up by email
  grafana user lookup admin@example.com`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.LookupUser(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:       %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Login:    %s\n", result.Login)
				fmt.Fprintf(f.IOStreams.Out, "Email:    %s\n", result.Email)
				fmt.Fprintf(f.IOStreams.Out, "Name:     %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Admin:    %v\n", result.IsAdmin)
				fmt.Fprintf(f.IOStreams.Out, "Disabled: %v\n", result.IsDisabled)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
