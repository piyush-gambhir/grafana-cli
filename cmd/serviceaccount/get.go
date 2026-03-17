package serviceaccount

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdServiceAccountGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a service account by ID",
		Long: `Retrieve a service account by its numeric ID.

Examples:
  # Get service account details
  grafana service-account get 10

  # Get as JSON
  grafana service-account get 10 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetServiceAccount(context.Background(), id)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:       %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Name:     %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Login:    %s\n", result.Login)
				fmt.Fprintf(f.IOStreams.Out, "Role:     %s\n", result.Role)
				fmt.Fprintf(f.IOStreams.Out, "Tokens:   %d\n", result.Tokens)
				fmt.Fprintf(f.IOStreams.Out, "Disabled: %v\n", result.IsDisabled)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
