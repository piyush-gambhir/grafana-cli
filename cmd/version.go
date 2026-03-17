package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/build"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "grafana-cli version %s\n", build.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", build.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", build.Date)
		},
	}
}
