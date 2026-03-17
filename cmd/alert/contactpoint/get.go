package contactpoint

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdContactPointGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <uid>",
		Short: "Get a contact point by UID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			cp, err := c.GetContactPoint(context.Background(), args[0])
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "UID:        %s\n", cp.UID)
				fmt.Fprintf(f.IOStreams.Out, "Name:       %s\n", cp.Name)
				fmt.Fprintf(f.IOStreams.Out, "Type:       %s\n", cp.Type)
				fmt.Fprintf(f.IOStreams.Out, "Provenance: %s\n", cp.Provenance)
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, cp, nil)
		},
	}
}
