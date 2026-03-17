package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func newCmdDashboardExport(f *cmdutil.Factory) *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "export <uid>",
		Short: "Export a dashboard JSON",
		Long:  "Export the full dashboard JSON to stdout or a file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetDashboardByUID(context.Background(), args[0])
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(result.Dashboard, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling dashboard: %w", err)
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, append(data, '\n'), 0o644); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				fmt.Fprintf(f.IOStreams.Out, "Dashboard exported to %s\n", outputFile)
				return nil
			}

			fmt.Fprintln(f.IOStreams.Out, string(data))
			return nil
		},
	}

	cmd.Flags().StringVar(&outputFile, "output-file", "", "Write output to file instead of stdout")

	return cmd
}
