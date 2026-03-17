package snapshot

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdSnapshotGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a snapshot by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetSnapshot(context.Background(), args[0])
			if err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
