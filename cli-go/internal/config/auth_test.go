package config

import (
	"testing"
)

func TestResolveAuth_FlagsOverrideAll(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://env-url")
	t.Setenv("GRAFANA_TOKEN", "env-token")

	profile := &Profile{
		URL:   "http://profile-url",
		Token: "profile-token",
	}

	rc := Resolve("http://flag-url", "flag-token", "", "", 0, profile, Defaults{})
	if rc.URL != "http://flag-url" {
		t.Errorf("got URL %q, want %q", rc.URL, "http://flag-url")
	}
	if rc.Token != "flag-token" {
		t.Errorf("got Token %q, want %q", rc.Token, "flag-token")
	}
}

func TestResolveAuth_EnvOverridesConfig(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://env-url")
	t.Setenv("GRAFANA_TOKEN", "env-token")
	t.Setenv("GRAFANA_USERNAME", "env-user")
	t.Setenv("GRAFANA_PASSWORD", "env-pass")
	t.Setenv("GRAFANA_ORG_ID", "99")

	profile := &Profile{
		URL:      "http://profile-url",
		Token:    "profile-token",
		Username: "profile-user",
		Password: "profile-pass",
		OrgID:    1,
	}

	rc := Resolve("", "", "", "", 0, profile, Defaults{})
	if rc.URL != "http://env-url" {
		t.Errorf("got URL %q, want %q", rc.URL, "http://env-url")
	}
	if rc.Token != "env-token" {
		t.Errorf("got Token %q, want %q", rc.Token, "env-token")
	}
	if rc.Username != "env-user" {
		t.Errorf("got Username %q, want %q", rc.Username, "env-user")
	}
	if rc.Password != "env-pass" {
		t.Errorf("got Password %q, want %q", rc.Password, "env-pass")
	}
	if rc.OrgID != 99 {
		t.Errorf("got OrgID %d, want %d", rc.OrgID, 99)
	}
}

func TestResolveAuth_ConfigFallback(t *testing.T) {
	// Clear any env vars.
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	t.Setenv("GRAFANA_USERNAME", "")
	t.Setenv("GRAFANA_PASSWORD", "")
	t.Setenv("GRAFANA_ORG_ID", "")

	profile := &Profile{
		URL:      "http://profile-url",
		Token:    "profile-token",
		Username: "profile-user",
		Password: "profile-pass",
		OrgID:    7,
	}

	rc := Resolve("", "", "", "", 0, profile, Defaults{Output: "yaml"})
	if rc.URL != "http://profile-url" {
		t.Errorf("got URL %q, want %q", rc.URL, "http://profile-url")
	}
	if rc.Token != "profile-token" {
		t.Errorf("got Token %q, want %q", rc.Token, "profile-token")
	}
	if rc.OrgID != 7 {
		t.Errorf("got OrgID %d, want %d", rc.OrgID, 7)
	}
	if rc.Output != "yaml" {
		t.Errorf("got Output %q, want %q", rc.Output, "yaml")
	}
}

func TestResolveAuth_MissingURL_Error(t *testing.T) {
	// When no URL is provided from any source, it should be empty.
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	t.Setenv("GRAFANA_USERNAME", "")
	t.Setenv("GRAFANA_PASSWORD", "")
	t.Setenv("GRAFANA_ORG_ID", "")

	rc := Resolve("", "", "", "", 0, nil, Defaults{})
	if rc.URL != "" {
		t.Errorf("expected empty URL, got %q", rc.URL)
	}
	// The error is raised by NewClient, not Resolve. Verify URL is empty.
}

func TestResolveAuth_BearerToken(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	t.Setenv("GRAFANA_USERNAME", "")
	t.Setenv("GRAFANA_PASSWORD", "")
	t.Setenv("GRAFANA_ORG_ID", "")

	rc := Resolve("http://localhost:3000", "glsa_test_token", "", "", 0, nil, Defaults{})
	if rc.Token != "glsa_test_token" {
		t.Errorf("got Token %q, want %q", rc.Token, "glsa_test_token")
	}
	if rc.Username != "" {
		t.Errorf("expected empty Username, got %q", rc.Username)
	}
}

func TestResolveAuth_BasicAuth(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	t.Setenv("GRAFANA_USERNAME", "")
	t.Setenv("GRAFANA_PASSWORD", "")
	t.Setenv("GRAFANA_ORG_ID", "")

	rc := Resolve("http://localhost:3000", "", "admin", "password123", 0, nil, Defaults{})
	if rc.Username != "admin" {
		t.Errorf("got Username %q, want %q", rc.Username, "admin")
	}
	if rc.Password != "password123" {
		t.Errorf("got Password %q, want %q", rc.Password, "password123")
	}
	if rc.Token != "" {
		t.Errorf("expected empty Token, got %q", rc.Token)
	}
}

func TestResolveAuth_OrgID(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	t.Setenv("GRAFANA_USERNAME", "")
	t.Setenv("GRAFANA_PASSWORD", "")
	t.Setenv("GRAFANA_ORG_ID", "")

	// Flag OrgID takes precedence.
	rc := Resolve("http://localhost:3000", "", "", "", 42, &Profile{OrgID: 10}, Defaults{})
	if rc.OrgID != 42 {
		t.Errorf("got OrgID %d, want %d", rc.OrgID, 42)
	}

	// Env OrgID is next.
	t.Setenv("GRAFANA_ORG_ID", "99")
	rc = Resolve("http://localhost:3000", "", "", "", 0, &Profile{OrgID: 10}, Defaults{})
	if rc.OrgID != 99 {
		t.Errorf("got OrgID %d, want %d", rc.OrgID, 99)
	}

	// Profile OrgID is last.
	t.Setenv("GRAFANA_ORG_ID", "")
	rc = Resolve("http://localhost:3000", "", "", "", 0, &Profile{OrgID: 10}, Defaults{})
	if rc.OrgID != 10 {
		t.Errorf("got OrgID %d, want %d", rc.OrgID, 10)
	}
}
