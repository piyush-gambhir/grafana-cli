package datasource

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDatasourceUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <id>",
		Short:       "Update a datasource",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update an existing datasource from a JSON or YAML file.

The file must contain the full datasource definition. Use
"grafana datasource get <uid> -o json" to retrieve the current state,
modify it, and pass it back.

Examples:
  # Update datasource with numeric ID 5
  grafana datasource update 5 -f updated-ds.json

  # Typical workflow: export, edit, update
  grafana datasource get P1234 -o json > ds.json
  # edit ds.json
  grafana datasource update 5 -f ds.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid datasource ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.DatasourceCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.UpdateDatasource(context.Background(), id, req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Datasource updated: %s (ID: %d)\n", result.Name, result.ID)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
