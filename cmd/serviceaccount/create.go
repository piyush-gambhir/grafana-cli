package serviceaccount

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdServiceAccountCreate(f *cmdutil.Factory) *cobra.Command {
	var file string
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create a service account",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new service account from a JSON or YAML file.

The file must contain name and role fields. Role can be Viewer, Editor,
or Admin.

Examples:
  # Create a service account
  grafana service-account create -f sa.json

  # Example JSON: {"name":"ci-bot","role":"Editor"}

  # Create idempotently (no error if already exists)
  grafana service-account create -f sa.json --if-not-exists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.ServiceAccountCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateServiceAccount(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: service account already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Service account created: %s (ID: %d)\n", result.Name, result.ID)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)
	cmdutil.AddIfNotExistsFlag(cmd, &ifNotExists)

	return cmd
}
