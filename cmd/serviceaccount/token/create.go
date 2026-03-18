package token

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTokenCreate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "create <service-account-id>",
		Short:       "Create a token for a service account",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new API token for a service account.

The token key is only shown once at creation time. Save it immediately.
The file must contain a "name" field and optionally "secondsToLive".

Examples:
  # Create a token
  grafana service-account token create 10 -f token.json

  # Example JSON: {"name":"deploy-token","secondsToLive":86400}`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			saID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.ServiceAccountTokenCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateServiceAccountToken(context.Background(), saID, req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Token created: %s (ID: %d)\n", result.Name, result.ID)
				if result.Key != "" {
					fmt.Fprintf(f.IOStreams.Out, "Key: %s\n", result.Key)
					fmt.Fprintln(f.IOStreams.Out, "NOTE: Save the key now. You will not be able to see it again.")
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
