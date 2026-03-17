package serviceaccount

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cmd/serviceaccount/token"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdServiceAccount returns the service account parent command.
func NewCmdServiceAccount(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-account",
		Short:   "Manage service accounts",
		Aliases: []string{"sa"},
	}

	cmd.AddCommand(newCmdServiceAccountList(f))
	cmd.AddCommand(newCmdServiceAccountGet(f))
	cmd.AddCommand(newCmdServiceAccountCreate(f))
	cmd.AddCommand(newCmdServiceAccountUpdate(f))
	cmd.AddCommand(newCmdServiceAccountDelete(f))
	cmd.AddCommand(token.NewCmdToken(f))

	return cmd
}
