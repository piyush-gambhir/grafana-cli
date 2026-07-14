package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Profiles == nil {
		t.Fatal("expected Profiles map to be initialized")
	}
	if cfg.Defaults.Output != "table" {
		t.Errorf("got default output %q, want %q", cfg.Defaults.Output, "table")
	}
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `current_profile: prod
profiles:
  prod:
    url: https://grafana.example.com
    token: glsa_xxxx
    org_id: 5
  staging:
    url: https://staging.grafana.example.com
    username: admin
    password: secret
defaults:
  output: json
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing test config: %v", err)
	}

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CurrentProfile != "prod" {
		t.Errorf("got current profile %q, want %q", cfg.CurrentProfile, "prod")
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("got %d profiles, want 2", len(cfg.Profiles))
	}
	prod := cfg.Profiles["prod"]
	if prod.URL != "https://grafana.example.com" {
		t.Errorf("got prod URL %q, want %q", prod.URL, "https://grafana.example.com")
	}
	if prod.Token != "glsa_xxxx" {
		t.Errorf("got prod token %q, want %q", prod.Token, "glsa_xxxx")
	}
	if prod.OrgID != 5 {
		t.Errorf("got prod org_id %d, want %d", prod.OrgID, 5)
	}
	if cfg.Defaults.Output != "json" {
		t.Errorf("got default output %q, want %q", cfg.Defaults.Output, "json")
	}
}

func TestSaveConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.yaml")

	cfg := &Config{
		CurrentProfile: "default",
		Profiles: map[string]Profile{
			"default": {URL: "http://localhost:3000", Token: "test-token"},
		},
		Defaults: Defaults{Output: "json"},
	}

	if err := cfg.SaveTo(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read it back.
	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error loading saved config: %v", err)
	}
	if loaded.CurrentProfile != "default" {
		t.Errorf("got current profile %q, want %q", loaded.CurrentProfile, "default")
	}
	if loaded.Profiles["default"].Token != "test-token" {
		t.Errorf("got token %q, want %q", loaded.Profiles["default"].Token, "test-token")
	}

	// Verify file permissions (owner read/write only).
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("got permissions %o, want 600", perm)
	}
}

func TestAddProfile(t *testing.T) {
	cfg := &Config{Profiles: make(map[string]Profile)}

	err := cfg.CreateProfile("prod", Profile{URL: "https://grafana.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Profiles["prod"]; !ok {
		t.Error("expected prod profile to exist")
	}

	// Adding duplicate should fail.
	err = cfg.CreateProfile("prod", Profile{URL: "https://other.example.com"})
	if err == nil {
		t.Error("expected error for duplicate profile")
	}
}

func TestGetProfile(t *testing.T) {
	cfg := &Config{
		CurrentProfile: "prod",
		Profiles: map[string]Profile{
			"prod": {URL: "https://grafana.example.com", Token: "token123"},
		},
	}

	p := cfg.CurrentProfileConfig()
	if p == nil {
		t.Fatal("expected profile to be found")
	}
	if p.URL != "https://grafana.example.com" {
		t.Errorf("got URL %q, want %q", p.URL, "https://grafana.example.com")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	cfg := &Config{
		CurrentProfile: "nonexistent",
		Profiles:       make(map[string]Profile),
	}

	p := cfg.CurrentProfileConfig()
	if p != nil {
		t.Error("expected nil profile for nonexistent name")
	}

	// Also test empty current profile.
	cfg.CurrentProfile = ""
	p = cfg.CurrentProfileConfig()
	if p != nil {
		t.Error("expected nil profile for empty current profile")
	}
}

func TestDeleteProfile(t *testing.T) {
	cfg := &Config{
		CurrentProfile: "prod",
		Profiles: map[string]Profile{
			"prod":    {URL: "https://grafana.example.com"},
			"staging": {URL: "https://staging.grafana.example.com"},
		},
	}

	// Delete current profile should unset it.
	err := cfg.DeleteProfile("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Profiles["prod"]; ok {
		t.Error("expected prod profile to be removed")
	}
	if cfg.CurrentProfile != "" {
		t.Errorf("expected empty current profile, got %q", cfg.CurrentProfile)
	}

	// Delete non-current profile.
	cfg.CurrentProfile = "staging"
	err = cfg.DeleteProfile("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Delete non-existent should error.
	err = cfg.DeleteProfile("nonexistent")
	if err == nil {
		t.Error("expected error for deleting nonexistent profile")
	}
}

func TestSetCurrentProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"prod":    {URL: "https://grafana.example.com"},
			"staging": {URL: "https://staging.grafana.example.com"},
		},
	}

	err := cfg.SetCurrentProfile("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CurrentProfile != "staging" {
		t.Errorf("got current profile %q, want %q", cfg.CurrentProfile, "staging")
	}

	// Setting nonexistent should error.
	err = cfg.SetCurrentProfile("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent profile")
	}
}

func TestConfigDir(t *testing.T) {
	// Test with XDG_CONFIG_HOME set.
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-test")
	dir := ConfigDir()
	expected := filepath.Join("/tmp/xdg-test", "grafana-cli")
	if dir != expected {
		t.Errorf("got config dir %q, want %q", dir, expected)
	}

	// Test with XDG_CONFIG_HOME unset (falls back to ~/.config).
	t.Setenv("XDG_CONFIG_HOME", "")
	dir = ConfigDir()
	home, _ := os.UserHomeDir()
	expected = filepath.Join(home, ".config", "grafana-cli")
	if dir != expected {
		t.Errorf("got config dir %q, want %q", dir, expected)
	}
}
