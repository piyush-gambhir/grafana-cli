package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/config"
)

func newLoginCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Interactively log in to a Grafana instance and save a profile",
		Long: `Interactively configure and test a connection to a Grafana instance.

Prompts for the server URL, authentication method (token or basic auth),
and credentials. Tests the connection, then saves the configuration as
a named profile.

Examples:
  # Start interactive login
  grafana login`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if f.NoInput {
				return fmt.Errorf("interactive input required but --no-input is set. Use environment variables (GRAFANA_URL, GRAFANA_TOKEN) instead of 'grafana login'.")
			}

			reader := bufio.NewReader(os.Stdin)
			out := f.IOStreams.Out

			// Prompt for URL.
			fmt.Fprint(out, "Grafana URL: ")
			urlStr, _ := reader.ReadString('\n')
			urlStr = strings.TrimSpace(urlStr)
			if urlStr == "" {
				return fmt.Errorf("URL is required")
			}

			// Prompt for auth method.
			fmt.Fprint(out, "Auth method (token/basic) [token]: ")
			authMethod, _ := reader.ReadString('\n')
			authMethod = strings.TrimSpace(authMethod)
			if authMethod == "" {
				authMethod = "token"
			}

			profile := config.Profile{URL: urlStr}

			switch authMethod {
			case "token":
				token, err := readLoginSecret(reader, out, "API Token: ")
				if err != nil {
					return err
				}
				if token == "" {
					return fmt.Errorf("token is required")
				}
				profile.Token = token
			case "basic":
				fmt.Fprint(out, "Username: ")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				password, err := readLoginSecret(reader, out, "Password: ")
				if err != nil {
					return err
				}
				if username == "" || password == "" {
					return fmt.Errorf("username and password are required")
				}
				profile.Username = username
				profile.Password = password
			default:
				return fmt.Errorf("invalid auth method: %s (use token or basic)", authMethod)
			}

			// Test the connection.
			fmt.Fprintln(out, "Testing connection...")
			resolved := &config.ResolvedConfig{
				URL:      profile.URL,
				Token:    profile.Token,
				Username: profile.Username,
				Password: profile.Password,
			}
			c, err := client.NewClient(resolved)
			if err != nil {
				return fmt.Errorf("creating client: %w", err)
			}

			resp, err := c.Get(cmd.Context(), "/api/org/")
			if err != nil {
				return fmt.Errorf("testing connection: %w", err)
			}
			var orgResult struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			if err := resp.JSON(&orgResult); err != nil {
				// Try health endpoint as fallback.
				resp2, err2 := c.Get(cmd.Context(), "/api/health")
				if err2 != nil {
					return fmt.Errorf("connection test failed: %w", err2)
				}
				if err2 := resp2.Error(); err2 != nil {
					return fmt.Errorf("connection test failed: %w", err2)
				}
				fmt.Fprintln(out, "Connection successful (health check passed)")
			} else {
				fmt.Fprintf(out, "Connection successful! Org: %s (ID: %d)\n", orgResult.Name, orgResult.ID)
			}

			// Prompt for profile name.
			fmt.Fprint(out, "Profile name [default]: ")
			profileName, _ := reader.ReadString('\n')
			profileName = strings.TrimSpace(profileName)
			if profileName == "" {
				profileName = "default"
			}

			if err := config.Update(func(cfg *config.Config) error {
				cfg.Profiles[profileName] = profile
				cfg.CurrentProfile = profileName
				return nil
			}); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(out, "Profile %q saved and set as current.\n", profileName)
			return nil
		},
	}
}

func readLoginSecret(reader *bufio.Reader, out io.Writer, label string) (string, error) {
	fmt.Fprint(out, label)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		secret, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(out)
		return strings.TrimSpace(string(secret)), err
	}
	secret, err := reader.ReadString('\n')
	return strings.TrimSpace(secret), err
}
