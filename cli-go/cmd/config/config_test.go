package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
)

func TestConfigViewMasksSecrets(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	if err := os.MkdirAll(filepath.Join(configHome, "grafana-cli"), 0700); err != nil {
		t.Fatal(err)
	}
	data := `current_profile: prod
profiles:
  prod:
    url: https://grafana.example.com
    token: glsa_secrettoken
  basic:
    url: https://grafana.example.com
    username: admin
    password: hunter2
`
	if err := os.WriteFile(filepath.Join(configHome, "grafana-cli", "config.yaml"), []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	f := &cmdutil.Factory{IOStreams: cmdutil.IOStreams{Out: &out, ErrOut: &out}}
	cmd := newCmdConfigView(f)
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.String()
	if strings.Contains(got, "glsa_secrettoken") || strings.Contains(got, "hunter2") {
		t.Errorf("config view printed a secret:\n%s", got)
	}
	if !strings.Contains(got, "token: '****'") || !strings.Contains(got, "password: '****'") {
		t.Errorf("config view did not mask secrets:\n%s", got)
	}
	if !strings.Contains(got, "username: admin") || !strings.Contains(got, "current_profile: prod") {
		t.Errorf("config view dropped non-secret fields:\n%s", got)
	}
}
