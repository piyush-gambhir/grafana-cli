package libraryelement

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdLibraryElementGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a library element by UID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetLibraryElement(context.Background(), args[0])
			if err != nil {
				return err
			}

			le := result.Result
			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "UID:         %s\n", le.UID)
				fmt.Fprintf(f.IOStreams.Out, "Name:        %s\n", le.Name)
				fmt.Fprintf(f.IOStreams.Out, "Type:        %s\n", le.Type)
				fmt.Fprintf(f.IOStreams.Out, "Kind:        %d\n", le.Kind)
				fmt.Fprintf(f.IOStreams.Out, "Description: %s\n", le.Description)
				fmt.Fprintf(f.IOStreams.Out, "Version:     %d\n", le.Version)
				fmt.Fprintf(f.IOStreams.Out, "Folder:      %s\n", le.Meta.FolderName)
				fmt.Fprintf(f.IOStreams.Out, "Connections: %d\n", le.Meta.ConnectedDashboards)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}
}
