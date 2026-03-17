package template

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTemplateGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get a notification template by name",
		Long: `Retrieve a notification template by its name.

In table mode, shows Name, Provenance, and the Template content.

Examples:
  # Get a template
  grafana alert template get "my-template"

  # Get as JSON
  grafana alert template get "my-template" -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAlertTemplate(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Name:       %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Provenance: %s\n", result.Provenance)
				fmt.Fprintf(f.IOStreams.Out, "Template:\n%s\n", result.Template)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
