package token

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdTokenList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <service-account-id>",
		Short:   "List tokens for a service account",
		Long: `List all API tokens for a service account.

The output includes ID, Name, Created date, Expiration date, and whether
the token has expired.

Examples:
  # List tokens for service account 10
  grafana service-account token list 10

  # Output as JSON
  grafana service-account token list 10 -o json`,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			saID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid service account ID: %s", args[0])
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListServiceAccountTokens(context.Background(), saID)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No tokens found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"ID", "Name", "Created", "Expiration", "Expired"},
				RowFunc: func(item interface{}) []string {
					t := item.(client.ServiceAccountToken)
					return []string{
						fmt.Sprintf("%d", t.ID),
						t.Name,
						t.Created,
						t.Expiration,
						fmt.Sprintf("%v", t.HasExpired),
					}
				},
			})
		},
	}
}
