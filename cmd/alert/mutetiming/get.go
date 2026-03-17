package mutetiming

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdMuteTimingGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get a mute timing by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetMuteTiming(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Name:       %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Intervals:  %d\n", len(result.TimeIntervals))
				fmt.Fprintf(f.IOStreams.Out, "Provenance: %s\n", result.Provenance)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
