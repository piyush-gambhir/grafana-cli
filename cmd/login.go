package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/config"
)

func newLoginCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Interactively log in to a Grafana instance and save a profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
				fmt.Fprint(out, "API Token: ")
				token, _ := reader.ReadString('\n')
				token = strings.TrimSpace(token)
				if token == "" {
					return fmt.Errorf("token is required")
				}
				profile.Token = token
			case "basic":
				fmt.Fprint(out, "Username: ")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				fmt.Fprint(out, "Password: ")
				password, _ := reader.ReadString('\n')
				password = strings.TrimSpace(password)
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

			resp, err := c.Get(context.Background(), "/api/org/")
			if err != nil {
				return fmt.Errorf("testing connection: %w", err)
			}
			var orgResult struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			if err := resp.JSON(&orgResult); err != nil {
				// Try health endpoint as fallback.
				resp2, err2 := c.Get(context.Background(), "/api/health")
				if err2 != nil {
					return fmt.Errorf("connection test failed: %w", err)
				}
				if err2 := resp2.Error(); err2 != nil {
					return fmt.Errorf("connection test failed: %w", err)
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

			// Save to config.
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Overwrite if exists.
			cfg.Profiles[profileName] = profile
			cfg.CurrentProfile = profileName

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(out, "Profile %q saved and set as current.\n", profileName)
			return nil
		},
	}
}
