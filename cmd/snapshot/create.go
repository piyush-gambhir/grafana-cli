package snapshot

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdSnapshotCreate(f *cmdutil.Factory) *cobra.Command {
	var file string
	var ifNotExists bool

	cmd := &cobra.Command{
		Use:         "create",
		Short:       "Create a snapshot",
		Annotations: map[string]string{"mutates": "true"},
		Long: `Create a dashboard snapshot from a JSON or YAML file.

The file must contain a "dashboard" object with the dashboard data.
Optionally include name, expires (seconds), and external flag.

Examples:
  # Create a snapshot
  grafana snapshot create -f snapshot.json

  # Create idempotently (no error if already exists)
  grafana snapshot create -f snapshot.json --if-not-exists`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file/-f is required")
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			var req client.SnapshotCreateRequest
			if err := cmdutil.UnmarshalInput(file, &req); err != nil {
				return err
			}

			result, err := c.CreateSnapshot(context.Background(), req)
			if err != nil {
				if ifNotExists && client.IsConflict(err) {
					if !f.Quiet {
						fmt.Fprintf(f.IOStreams.ErrOut, "Warning: snapshot already exists, skipping.\n")
					}
					return nil
				}
				return err
			}

			if f.Resolved.Output == "table" {
				if !f.Quiet {
					fmt.Fprintf(f.IOStreams.Out, "Snapshot created.\n")
					fmt.Fprintf(f.IOStreams.Out, "  Key:        %s\n", result.Key)
					fmt.Fprintf(f.IOStreams.Out, "  URL:        %s\n", result.URL)
					fmt.Fprintf(f.IOStreams.Out, "  Delete Key: %s\n", result.DeleteKey)
				}
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result, nil)
		},
	}

	cmdutil.AddFileFlag(cmd, &file)
	cmdutil.AddIfNotExistsFlag(cmd, &ifNotExists)

	return cmd
}
