package serviceaccount

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdServiceAccountUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <id>",
		Short:       "Update a service account",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update a service account from a JSON or YAML file.

Examples:
  # Update service account 10
  grafana service-account update 10 -f sa.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.ServiceAccountUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.UpdateServiceAccount(context.Background(), id, req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Service account updated: %s (ID: %d)\n", result.Name, result.ID)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
