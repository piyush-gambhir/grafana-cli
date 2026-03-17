package token

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdToken returns the service account token parent command.
func NewCmdToken(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage service account tokens",
	}

	cmd.AddCommand(newCmdTokenList(f))
	cmd.AddCommand(newCmdTokenCreate(f))
	cmd.AddCommand(newCmdTokenDelete(f))

	return cmd
}
