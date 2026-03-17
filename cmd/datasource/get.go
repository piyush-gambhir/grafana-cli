package datasource

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDatasourceGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a datasource by UID",
		Long: `Retrieve a single datasource by its UID.

In table mode, displays ID, UID, Name, Type, URL, Access mode, Default
status, and ReadOnly status. In JSON/YAML mode, returns the full datasource
object including jsonData configuration.

Examples:
  # Get datasource details
  grafana datasource get P1234

  # Get full datasource JSON (useful for creating update payloads)
  grafana datasource get P1234 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetDatasourceByUID(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:        %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "UID:       %s\n", result.UID)
				fmt.Fprintf(f.IOStreams.Out, "Name:      %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Type:      %s\n", result.Type)
				fmt.Fprintf(f.IOStreams.Out, "URL:       %s\n", result.URL)
				fmt.Fprintf(f.IOStreams.Out, "Access:    %s\n", result.Access)
				fmt.Fprintf(f.IOStreams.Out, "Default:   %v\n", result.IsDefault)
				fmt.Fprintf(f.IOStreams.Out, "ReadOnly:  %v\n", result.ReadOnly)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	return cmd
}
