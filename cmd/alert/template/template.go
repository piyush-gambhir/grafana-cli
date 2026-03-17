package template

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdTemplate returns the notification template parent command.
func NewCmdTemplate(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Manage notification templates",
		Long: `List, get, update, and delete notification templates.

Templates customize the content of alert notifications using Go templating
syntax. They are referenced by contact points.`,
		Aliases: []string{"tmpl"},
	}

	cmd.AddCommand(newCmdTemplateList(f))
	cmd.AddCommand(newCmdTemplateGet(f))
	cmd.AddCommand(newCmdTemplateUpdate(f))
	cmd.AddCommand(newCmdTemplateDelete(f))

	return cmd
}
