package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long: `Print the grafana-cli version, commit hash, and build date.

When a GitHub release check from the last 24 hours is cached, it also prints
the latest release and whether an update is available. It never contacts
GitHub itself; run "grafana update --check" for a fresh answer.

Examples:
  grafana version`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "grafana-cli version %s\n", build.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", build.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", build.Date)
			if !update.IsReleaseVersion(build.Version) {
				return
			}
			if info := update.Cached(build.Version, config.ConfigDir()); info != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  latest: %s\n", info.LatestVersion)
				fmt.Fprintf(cmd.OutOrStdout(), "  update_available: %t\n", info.Available)
			}
		},
	}
}
