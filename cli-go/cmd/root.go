package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/admin"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/alert"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/annotation"
	cmdconfig "github.com/piyush-gambhir/grafana-cli/cli-go/cmd/config"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/correlation"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/dashboard"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/datasource"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/folder"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/libraryelement"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/org"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/playlist"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/preferences"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/serviceaccount"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/snapshot"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/team"
	"github.com/piyush-gambhir/grafana-cli/cli-go/cmd/user"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
)

var (
	flagOutput   string
	flagProfile  string
	flagURL      string
	flagToken    string
	flagUsername string
	flagPassword string
	flagOrgID    int64
	flagReadOnly bool
	flagNoInput  bool
	flagQuiet    bool
	flagVerbose  bool
)

// OutputFormat is set during PersistentPreRunE and exported for use by main.go.
var OutputFormat string

// updateNoticeWait is how long PersistentPostRun waits for a background check
// this run started; the check itself times out after 3 seconds.
const updateNoticeWait = time.Second

// Test seams for the background update check.
var (
	stderrIsTerminal    = func() bool { return term.IsTerminal(int(os.Stderr.Fd())) }
	startUpdateCheck    = update.Start
	detectInstallMethod = func() string {
		exe, err := currentExecutable()
		if err != nil {
			return update.MethodSelf
		}
		return update.DetectInstallMethod(exe)
	}
)

// Execute is the main entry point for the CLI.
func Execute() error {
	if goos == "windows" {
		if exe, err := currentExecutable(); err == nil {
			cleanupOldExecutable(goos, exe)
		}
	}
	return newRootCmd().Execute()
}

// isTopLevel reports whether cmd is a direct child of the root named one of
// names, so `grafana update` matches but `grafana dashboard update` does not.
func isTopLevel(cmd *cobra.Command, names ...string) bool {
	if !cmd.HasParent() || cmd.Parent().HasParent() {
		return false
	}
	for _, name := range names {
		if cmd.Name() == name {
			return true
		}
	}
	return false
}

// skipsUpdateCheck lists the commands that never run the background check.
func skipsUpdateCheck(cmd *cobra.Command) bool {
	return strings.HasPrefix(cmd.Name(), "__complete") || isTopLevel(cmd, "update", "version", "completion", "help")
}

// loadAndResolveConfig loads the config file and resolves auth from flags/env/config.
func loadAndResolveConfig(cmd *cobra.Command) (*config.ResolvedConfig, *config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("loading config: %w", err)
	}

	// Determine which profile to use.
	profileName := flagProfile
	if profileName == "" {
		profileName = cfg.CurrentProfile
	}
	var profile *config.Profile
	if profileName != "" {
		p, ok := cfg.Profiles[profileName]
		if !ok && cmd.Name() != "login" {
			return nil, nil, fmt.Errorf("profile %q not found", profileName)
		}
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

	return resolved, cfg, nil
}

func envFlagEnabled(name string) bool {
	v := strings.TrimSpace(os.Getenv(name))
	return strings.EqualFold(v, "true") || v == "1"
}

// createClient sets up the HTTP client factory on the factory.
func createClient(f *cmdutil.Factory, resolved *config.ResolvedConfig) {
	f.Client = func() (*client.Client, error) {
		c, err := client.NewClient(resolved)
		if err != nil {
			return nil, err
		}
		if flagVerbose {
			c.EnableVerboseLogging(f.IOStreams.ErrOut)
		}
		return c, nil
	}
}

// checkPermissions enforces read-only and no-input checks.
func checkPermissions(cmd *cobra.Command, resolved *config.ResolvedConfig) error {
	effectiveReadOnly := resolved.ReadOnly // from env > config
	if flagReadOnly {
		effectiveReadOnly = true
	}
	if effectiveReadOnly && cmd.Annotations != nil && cmd.Annotations["mutates"] == "true" {
		return fmt.Errorf("command '%s' is blocked in read-only mode; remove read_only from the profile or disable the read-only environment setting to permit writes", cmd.CommandPath())
	}
	return nil
}

func newRootCmd() *cobra.Command {
	f := &cmdutil.Factory{
		IOStreams: cmdutil.DefaultIOStreams(),
	}

	// The background update check's answer, passed from PersistentPreRunE to PersistentPostRun.
	var updateResult <-chan *update.UpdateInfo

	rootCmd := &cobra.Command{
		Use:   "grafana",
		Short: "Grafana CLI - manage Grafana from the command line",
		Long: `A command-line interface for managing Grafana instances, dashboards, datasources, alerts, and more.

Full command reference (for agents/LLMs): https://github.com/piyush-gambhir/grafana-cli/blob/main/docs/llms.txt
Claude Code skill: https://github.com/piyush-gambhir/grafana-cli/blob/main/grafana/SKILL.md`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Check env vars for --no-input, --quiet, --verbose.
			if envFlagEnabled("GRAFANA_NO_INPUT") {
				flagNoInput = true
			}
			if !cmd.Flags().Changed("quiet") {
				flagQuiet = envFlagEnabled("GRAFANA_QUIET")
			}
			if !cmd.Flags().Changed("verbose") {
				flagVerbose = envFlagEnabled("GRAFANA_VERBOSE")
			}
			f.NoInput = flagNoInput
			f.Quiet = flagQuiet
			f.Verbose = flagVerbose

			// Start the background update check unless the command, the
			// environment, or the build rules it out (no network then).
			if !skipsUpdateCheck(cmd) && !update.NotifierDisabled(os.Getenv, build.Version, flagQuiet, stderrIsTerminal()) {
				updateResult = startUpdateCheck(build.Version, config.ConfigDir())
			}

			// Skip auth setup for top-level commands that don't need it.
			if isTopLevel(cmd, "version", "completion", "help", "update") {
				return nil
			}
			// Also skip for config subcommands.
			if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
				return nil
			}

			resolved, cfg, err := loadAndResolveConfig(cmd)
			if err != nil {
				return err
			}

			// Set exported OutputFormat for use by main.go error handler.
			OutputFormat = resolved.Output

			f.Resolved = resolved

			f.Config = func() (*config.Config, error) {
				return cfg, nil
			}

			createClient(f, resolved)

			return checkPermissions(cmd, resolved)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if updateResult == nil {
				return
			}
			// A cached answer is already in the channel, so this never waits
			// then. A check this run started (at most once a day) gets up to
			// updateNoticeWait: it is recorded before the request, so a fast
			// command that exited first would lose the day's notice.
			select {
			case info := <-updateResult:
				if info != nil && info.Available {
					update.Notify(cmd.ErrOrStderr(), info, config.ConfigDir(), detectInstallMethod())
				}
			case <-time.After(updateNoticeWait):
			}
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
	rootCmd.PersistentFlags().BoolVar(&flagReadOnly, "read-only", false, "Block write operations (safety mode for agents)")
	rootCmd.PersistentFlags().BoolVar(&flagNoInput, "no-input", false, "Disable all interactive prompts (for CI/agent use)")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress informational output")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose HTTP logging")

	// Register subcommands.
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUpdateCmd())
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
