package folder

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdFolderUpdate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:         "update <uid>",
		Short:       "Update a folder",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Update an existing folder's title or other properties.

The file must contain the updated folder properties (e.g. title).

Examples:
  # Update a folder
  grafana folder update folderUid123 -f updated-folder.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.FolderUpdateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.UpdateFolder(context.Background(), args[0], req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Folder updated: %s (UID: %s)\n", result.Title, result.UID)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
