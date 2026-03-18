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
	var (
		tag   string
		limit int64
	)

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List all annotation tags",
		Long: `List all unique annotation tags with their usage counts.

Examples:
  # List all tags
  grafana annotation tags

  # Filter by tag prefix
  grafana annotation tags --tag deploy

  # Limit results
  grafana annotation tags --limit 10

  # Output as JSON
  grafana annotation tags -o json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			result, err := c.GetAnnotationTags(context.Background(), tag, limit)
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

	cmd.Flags().StringVar(&tag, "tag", "", "Filter tags by prefix")
	cmd.Flags().Int64Var(&limit, "limit", 0, "Maximum number of tags to return")

	return cmd
}
