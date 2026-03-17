package folder

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdFolderCreate(f *cmdutil.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a folder",
		Long:  "Create a new folder from a JSON/YAML file or with flags.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.FolderCreateRequest
			if file != "" {
				if err := cmdutil.UnmarshalInput(file, &req); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("--file/-f is required")
			}

			result, err := c.CreateFolder(context.Background(), req)
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Folder created: %s (UID: %s)\n", result.Title, result.UID)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)

	return cmd
}
