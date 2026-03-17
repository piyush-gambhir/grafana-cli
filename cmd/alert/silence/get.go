package silence

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdSilenceGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a silence by ID",
		Long: `Retrieve a single silence by its ID.

Examples:
  # Get silence details
  grafana alert silence get silenceId123

  # Get as JSON
  grafana alert silence get silenceId123 -o json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetSilence(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:         %s\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "State:      %s\n", result.Status.State)
				fmt.Fprintf(f.IOStreams.Out, "Comment:    %s\n", result.Comment)
				fmt.Fprintf(f.IOStreams.Out, "Created By: %s\n", result.CreatedBy)
				fmt.Fprintf(f.IOStreams.Out, "Starts At:  %s\n", result.StartsAt)
				fmt.Fprintf(f.IOStreams.Out, "Ends At:    %s\n", result.EndsAt)
				fmt.Fprintf(f.IOStreams.Out, "Matchers:   %d\n", len(result.Matchers))
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
