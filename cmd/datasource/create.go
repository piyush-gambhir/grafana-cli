package datasource

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdDatasourceCreate(f *cmdutil.Factory) *cobra.Command {
	var file string
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create a datasource",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new datasource from a JSON or YAML file.

The file must contain at minimum: name, type, access, and url fields.
Sensitive fields like passwords should be placed in the secureJsonData
object.

Examples:
  # Create a Prometheus datasource
  grafana datasource create -f prometheus.json

  # Create from YAML
  grafana datasource create -f datasource.yaml

  # Read from stdin
  echo '{"name":"test","type":"prometheus","access":"proxy","url":"http://prom:9090"}' | grafana datasource create -f -

  # Create idempotently (no error if already exists)
  grafana datasource create -f prometheus.json --if-not-exists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.DatasourceCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateDatasource(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: datasource already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Datasource created: %s (ID: %d)\n", result.Name, result.ID)
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
