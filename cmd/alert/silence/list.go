package silence

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdSilenceList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List silences",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListSilences(context.Background())
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No silences found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "State", "Matchers", "Starts At", "Ends At", "Comment", "Created By"},
				RowFunc: func(item interface{}) []string {
					s := item.(client.Silence)
					var matchers []string
					for _, m := range s.Matchers {
						op := "="
						if !m.IsEqual {
							op = "!="
						}
						if m.IsRegex {
							op += "~"
						}
						matchers = append(matchers, fmt.Sprintf("%s%s%s", m.Name, op, m.Value))
					}
					return []string{s.ID, s.Status.State, strings.Join(matchers, ", "), s.StartsAt, s.EndsAt, s.Comment, s.CreatedBy}
				},
			})
		},
	}
}
