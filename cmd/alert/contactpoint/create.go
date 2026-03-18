package contactpoint

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdContactPointCreate(f *cmdutil.Factory) *cobra.Command {
	var file string
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create a contact point",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new contact point from a JSON or YAML file.

The file must contain name, type, and settings fields. The type determines
which settings are required (e.g. "email" needs "addresses", "slack"
needs "url").

Examples:
  # Create a contact point
  grafana alert contact-point create -f contact-point.json

  # Create idempotently (no error if already exists)
  grafana alert contact-point create -f contact-point.json --if-not-exists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.ContactPoint
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateContactPoint(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: contact point already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Contact point created: %s (UID: %s)\n", result.Name, result.UID)
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
