package annotation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdAnnotationUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an annotation",
		Long: `Update an annotation's text, tags, or time range.

Examples:
  # Update annotation 42
  grafana annotation update 42 -f annotation.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid annotation ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.AnnotationUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			if err := c.UpdateAnnotation(context.Background(), id, req); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "Annotation %d updated.\n", id)
			return nil
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
