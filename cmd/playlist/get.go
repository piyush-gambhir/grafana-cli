package playlist

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdPlaylistGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a playlist by UID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetPlaylist(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "UID:      %s\n", result.UID)
				fmt.Fprintf(f.IOStreams.Out, "Name:     %s\n", result.Name)
				fmt.Fprintf(f.IOStreams.Out, "Interval: %s\n", result.Interval)
				fmt.Fprintf(f.IOStreams.Out, "Items:    %d\n", len(result.Items))
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
