package folder

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdFolderGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a folder by UID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetFolder(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "ID:        %d\n", result.ID)
				fmt.Fprintf(f.IOStreams.Out, "UID:       %s\n", result.UID)
				fmt.Fprintf(f.IOStreams.Out, "Title:     %s\n", result.Title)
				fmt.Fprintf(f.IOStreams.Out, "URL:       %s\n", result.URL)
				fmt.Fprintf(f.IOStreams.Out, "Version:   %d\n", result.Version)
				fmt.Fprintf(f.IOStreams.Out, "Created:   %s\n", result.Created)
				fmt.Fprintf(f.IOStreams.Out, "Updated:   %s\n", result.Updated)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
