package contactpoint

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdContactPointUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <uid>",
		Short:       "Update a contact point",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update an existing contact point from a JSON or YAML file.

Examples:
  # Update a contact point
  grafana alert contact-point update cpUid123 -f updated-cp.json`,
		Args:  cobra.ExactArgs(1),
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

			if err := c.UpdateContactPoint(context.Background(), args[0], req); err != nil {
				return err
			}

			if !f.Quiet {
				fmt.Fprintf(f.IOStreams.Out, "Contact point %q updated.\n", args[0])
			}
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
