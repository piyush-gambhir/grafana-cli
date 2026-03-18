package annotation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAnnotationCreate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create an annotation",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a new annotation from a JSON or YAML file.

The file must contain a "text" field. Optionally include dashboardId,
panelId, time (epoch ms), timeEnd, and tags.

Examples:
  # Create an annotation
  grafana annotation create -f annotation.json

  # Example JSON: {"text":"Deployed v1.2.3","tags":["deploy"]}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.AnnotationCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateAnnotation(context.Background(), req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Annotation created (ID: %d): %s\n", result.ID, result.Message)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
