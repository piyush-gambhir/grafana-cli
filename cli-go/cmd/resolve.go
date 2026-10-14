package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/output"
)

// resolution describes what arguments would run, using the real command parser.
type resolution struct {
	Command     string            `json:"command"`
	Effect      string            `json:"effect"`
	Conditional bool              `json:"conditional"`
	Interactive bool              `json:"interactive"`
	Secret      bool              `json:"secret_output"`
	Flags       map[string]string `json:"flags"`
	Args        []string          `json:"args"`
}

func newCommandsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "commands", Short: "Offline command discovery for agents"}
	cmd.AddCommand(&cobra.Command{
		Use:                "resolve -- <arguments>",
		Short:              "Show the command, effect and flags that arguments would run, without running them",
		DisableFlagParsing: true,
		Example:            "  grafana commands resolve -- dashboard list\n  grafana commands resolve -- -o json dashboard delete abc123 --confirm",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bypass config, credentials and the background update check.
			OutputFormat = "json"
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			for i, arg := range args {
				if arg == "--" {
					args = args[i+1:]
					break
				}
			}
			probe := newRootCmd()
			probe.SetIn(strings.NewReader(""))
			probe.SetOut(io.Discard)
			probe.SetErr(io.Discard)
			target, rest, err := probe.Find(args)
			if err != nil || target == probe || !target.Runnable() {
				return fmt.Errorf("arguments do not name a runnable grafana command")
			}
			if err := target.ParseFlags(rest); err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), OutputFormat, commandResolution(target), nil)
		},
	})
	return cmd
}

func commandResolution(target *cobra.Command) resolution {
	res := resolution{
		Command: target.CommandPath(), Effect: "read",
		Flags: map[string]string{}, Args: append([]string{}, target.Flags().Args()...),
	}
	ann := target.Annotations
	flagSet := func(name string) bool { return name != "" && target.Flags().Changed(name) }
	switch {
	case ann["mutates"] == "true" || ann["mutates"] == "conditional":
		res.Effect, res.Conditional = "remote_write", ann["mutates"] == "conditional"
	case flagSet(ann["mutates-when-flag"]):
		res.Effect = "remote_write"
	case ann["writes-local"] == "true" || ann["writes-local"] == "conditional" || ann["self-update"] != "" || flagSet(ann["writes-local-when-flag"]):
		res.Effect, res.Conditional = "local_write", ann["writes-local"] == "conditional"
	}
	res.Interactive = ann["interactive"] == "true"
	res.Secret = ann["secret-output"] == "true"
	target.Flags().Visit(func(f *pflag.Flag) { res.Flags[f.Name] = f.Value.String() })
	return res
}
