package annotation

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdAnnotation returns the annotation parent command.
func NewCmdAnnotation(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotation",
		Short: "Manage annotations",
		Long: `Create, list, get, update, and delete annotations, and view annotation tags.

Annotations mark points in time on Grafana graphs, commonly used for
deployments, incidents, or other significant events.`,
	}

	cmd.AddCommand(newCmdAnnotationList(f))
	cmd.AddCommand(newCmdAnnotationGet(f))
	cmd.AddCommand(newCmdAnnotationCreate(f))
	cmd.AddCommand(newCmdAnnotationUpdate(f))
	cmd.AddCommand(newCmdAnnotationDelete(f))
	cmd.AddCommand(newCmdAnnotationTags(f))

	return cmd
}
