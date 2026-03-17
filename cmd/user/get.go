package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdUserGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a user by ID",
		Long: `Retrieve a user by their numeric ID (requires server admin).

Examples:
  # Get user details
  grafana user get 5

  # Get as JSON
  grafana user get 5 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetUser(context.Background(), id)
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
