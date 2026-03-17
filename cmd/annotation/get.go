package annotation

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAnnotationGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an annotation by ID",
		Long: `Retrieve a single annotation by its numeric ID.

Examples:
  # Get annotation 42
  grafana annotation get 42

  # Get as JSON
  grafana annotation get 42 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid annotation ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAnnotation(context.Background(), id)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:           %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "Dashboard ID: %d\n", result.DashboardID)
				fmt.Fprintf(f.IOStreams.Out, "Panel ID:     %d\n", result.PanelID)
				fmt.Fprintf(f.IOStreams.Out, "Text:         %s\n", result.Text)
				fmt.Fprintf(f.IOStreams.Out, "Tags:         %s\n", strings.Join(result.Tags, ", "))
				fmt.Fprintf(f.IOStreams.Out, "Time:         %d\n", result.Time)
				fmt.Fprintf(f.IOStreams.Out, "Time End:     %d\n", result.TimeEnd)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
