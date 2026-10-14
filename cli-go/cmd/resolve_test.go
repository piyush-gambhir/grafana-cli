package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
	"github.com/spf13/cobra"
)

func TestCommandsResolve(t *testing.T) {
	// Invalid config and an enabled notifier must never affect offline resolution.
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	if err := os.MkdirAll(filepath.Join(configHome, "grafana-cli"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configHome, "grafana-cli", "config.yaml"), []byte("[invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	oldStart, oldTerminal := startUpdateCheck, stderrIsTerminal
	t.Cleanup(func() { startUpdateCheck, stderrIsTerminal = oldStart, oldTerminal })
	stderrIsTerminal = func() bool { return true }
	startUpdateCheck = func(string, string) <-chan *update.UpdateInfo {
		t.Fatal("resolve started an update check")
		return nil
	}
	for _, tc := range []struct {
		name string
		args []string
		want resolution
	}{
		{"read", []string{"dashboard", "list"}, resolution{Command: "grafana dashboard list", Effect: "read", Flags: map[string]string{}, Args: []string{}}},
		{"write with leading flag", []string{"-o", "json", "dashboard", "delete", "abc123", "--confirm"}, resolution{Command: "grafana dashboard delete", Effect: "remote_write", Flags: map[string]string{"output": "json", "confirm": "true"}, Args: []string{"abc123"}}},
		{"secret", []string{"sa", "token", "create", "10", "-f", "missing.json"}, resolution{Command: "grafana service-account token create", Effect: "remote_write", Secret: true, Flags: map[string]string{"file": "missing.json"}, Args: []string{"10"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			root := newRootCmd()
			root.SetOut(&out)
			root.SetArgs(append([]string{"commands", "resolve", "--"}, tc.args...))
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			var got resolution
			if err := json.Unmarshal(out.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON: %v: %s", err, &out)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
	for _, args := range [][]string{{"nonsense"}, {}, {"dashboard"}, {"dashboard", "list", "--unknown"}, {"dashboard", "list", "--limit", "invalid"}} {
		root := newRootCmd()
		root.SetArgs(append([]string{"commands", "resolve", "--"}, args...))
		if err := root.Execute(); err == nil {
			t.Errorf("%v: expected an error", args)
		}
		if OutputFormat != "json" {
			t.Errorf("%v: error format = %q, want json", args, OutputFormat)
		}
	}
}

func TestCommandsResolveAnnotations(t *testing.T) {
	for _, tc := range []struct {
		name        string
		annotations map[string]string
		args        []string
		effect      string
		conditional bool
	}{
		{"read", nil, nil, "read", false},
		{"remote conditional", map[string]string{"mutates": "conditional", "writes-local": "true"}, nil, "remote_write", true},
		{"remote flag absent", map[string]string{"mutates-when-flag": "write"}, nil, "read", false},
		{"remote flag explicitly false", map[string]string{"mutates-when-flag": "write"}, []string{"--write=false"}, "remote_write", false},
		{"local", map[string]string{"writes-local": "true"}, nil, "local_write", false},
		{"local conditional", map[string]string{"writes-local": "conditional"}, nil, "local_write", true},
		{"self update", map[string]string{"self-update": "true"}, nil, "local_write", false},
		{"local flag absent", map[string]string{"writes-local-when-flag": "write"}, nil, "read", false},
		{"local flag", map[string]string{"writes-local-when-flag": "write"}, []string{"--write"}, "local_write", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "probe", Annotations: tc.annotations}
			cmd.Flags().Bool("write", false, "Write output")
			if err := cmd.ParseFlags(tc.args); err != nil {
				t.Fatal(err)
			}
			got := commandResolution(cmd)
			if got.Effect != tc.effect || got.Conditional != tc.conditional {
				t.Errorf("got %+v, want effect=%s conditional=%t", got, tc.effect, tc.conditional)
			}
		})
	}
}
