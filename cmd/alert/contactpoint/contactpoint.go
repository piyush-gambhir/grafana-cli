package contactpoint

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdContactPoint returns the contact point parent command.
func NewCmdContactPoint(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contact-point",
		Short:   "Manage contact points",
		Long: `Create, list, get, update, and delete alert notification contact points.

Contact points define where alert notifications are sent (email, Slack,
PagerDuty, webhook, etc.).`,
		Aliases: []string{"cp"},
	}

	cmd.AddCommand(newCmdContactPointList(f))
	cmd.AddCommand(newCmdContactPointGet(f))
	cmd.AddCommand(newCmdContactPointCreate(f))
	cmd.AddCommand(newCmdContactPointUpdate(f))
	cmd.AddCommand(newCmdContactPointDelete(f))

	return cmd
}
