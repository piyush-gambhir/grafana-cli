package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cmd/admin"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert"
	"github.com/piyush-gambhir/grafana-cli/cmd/annotation"
	cmdconfig "github.com/piyush-gambhir/grafana-cli/cmd/config"
	"github.com/piyush-gambhir/grafana-cli/cmd/correlation"
	"github.com/piyush-gambhir/grafana-cli/cmd/dashboard"
	"github.com/piyush-gambhir/grafana-cli/cmd/datasource"
	"github.com/piyush-gambhir/grafana-cli/cmd/folder"
	"github.com/piyush-gambhir/grafana-cli/cmd/libraryelement"
	"github.com/piyush-gambhir/grafana-cli/cmd/org"
	"github.com/piyush-gambhir/grafana-cli/cmd/playlist"
	"github.com/piyush-gambhir/grafana-cli/cmd/preferences"
	"github.com/piyush-gambhir/grafana-cli/cmd/serviceaccount"
	"github.com/piyush-gambhir/grafana-cli/cmd/snapshot"
	"github.com/piyush-gambhir/grafana-cli/cmd/team"
	"github.com/piyush-gambhir/grafana-cli/cmd/user"
	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/config"
)

var (
	flagOutput   string
	flagProfile  string
	flagURL      string
	flagToken    string
	flagUsername string
	flagPassword string
	flagOrgID    int64
)

// Execute is the main entry point for the CLI.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	f := &cmdutil.Factory{
		IOStreams: cmdutil.DefaultIOStreams(),
	}

	rootCmd := &cobra.Command{
		Use:   "grafana",
		Short: "Grafana CLI - manage Grafana from the command line",
		Long:  "A command-line interface for managing Grafana instances, dashboards, datasources, alerts, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip auth setup for commands that don't need it.
			if cmd.Name() == "version" || cmd.Name() == "completion" || cmd.Name() == "help" {
				return nil
			}
			// Also skip for config subcommands.
			if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
				return nil
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Determine which profile to use.
			profileName := flagProfile
			if profileName == "" {
				profileName = cfg.CurrentProfile
			}
			var profile *config.Profile
			if profileName != "" {
				p, ok := cfg.Profiles[profileName]
				if ok {
					profile = &p
				}
			}

			// Determine output format.
			output := flagOutput
			if output == "" {
				output = cfg.Defaults.Output
			}

			// Resolve configuration.
			resolved := config.Resolve(flagURL, flagToken, flagUsername, flagPassword, flagOrgID, profile, cfg.Defaults)
			if output != "" {
				resolved.Output = output
			}
			f.Resolved = resolved

			f.Config = func() (*config.Config, error) {
				return cfg, nil
			}

			f.Client = func() (*client.Client, error) {
				return client.NewClient(resolved)
			}

			return nil
		},
	}

	// Global persistent flags.
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "", "Output format: table, json, yaml")
	rootCmd.PersistentFlags().StringVar(&flagProfile, "profile", "", "Configuration profile to use")
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Grafana server URL")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token or service account token")
	rootCmd.PersistentFlags().StringVar(&flagUsername, "username", "", "Username for basic auth")
	rootCmd.PersistentFlags().StringVar(&flagPassword, "password", "", "Password for basic auth")
	rootCmd.PersistentFlags().Int64Var(&flagOrgID, "org-id", 0, "Organization ID")

	// Register subcommands.
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newLoginCmd(f))
	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(cmdconfig.NewCmdConfig(f))
	rootCmd.AddCommand(dashboard.NewCmdDashboard(f))
	rootCmd.AddCommand(datasource.NewCmdDatasource(f))
	rootCmd.AddCommand(folder.NewCmdFolder(f))
	rootCmd.AddCommand(alert.NewCmdAlert(f))
	rootCmd.AddCommand(org.NewCmdOrg(f))
	rootCmd.AddCommand(team.NewCmdTeam(f))
	rootCmd.AddCommand(user.NewCmdUser(f))
	rootCmd.AddCommand(serviceaccount.NewCmdServiceAccount(f))
	rootCmd.AddCommand(annotation.NewCmdAnnotation(f))
	rootCmd.AddCommand(snapshot.NewCmdSnapshot(f))
	rootCmd.AddCommand(playlist.NewCmdPlaylist(f))
	rootCmd.AddCommand(libraryelement.NewCmdLibraryElement(f))
	rootCmd.AddCommand(correlation.NewCmdCorrelation(f))
	rootCmd.AddCommand(admin.NewCmdAdmin(f))
	rootCmd.AddCommand(preferences.NewCmdPreferences(f))

	return rootCmd
}
