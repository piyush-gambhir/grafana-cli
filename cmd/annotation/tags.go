package annotation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdAnnotationTags(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "List all annotation tags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAnnotationTags(context.Background())
			if err != nil {
				return err
			}

			if len(result.Result) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No annotation tags found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, result.Result, &output.TableDef{
				Headers: []string{"Tag", "Count"},
				RowFunc: func(item interface{}) []string {
					t := item.(client.AnnotationTag)
					return []string{t.Tag, fmt.Sprintf("%d", t.Count)}
				},
			})
		},
	}
}
