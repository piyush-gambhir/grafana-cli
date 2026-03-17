package correlation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdCorrelationGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <source-uid> <correlation-uid>",
		Short: "Get a correlation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetCorrelation(context.Background(), args[0], args[1])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "UID:         %s\n", result.UID)
				fmt.Fprintf(f.IOStreams.Out, "Source UID:  %s\n", result.SourceUID)
				fmt.Fprintf(f.IOStreams.Out, "Target UID:  %s\n", result.TargetUID)
				fmt.Fprintf(f.IOStreams.Out, "Label:       %s\n", result.Label)
				fmt.Fprintf(f.IOStreams.Out, "Description: %s\n", result.Description)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
