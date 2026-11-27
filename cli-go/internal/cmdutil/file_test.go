package cmdutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
)

func TestUnmarshalInput_YAMLUsesJSONFieldNames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "policy.yaml")
	data := `receiver: default
group_by: [alertname]
routes:
  - receiver: team-a
    object_matchers: [[team, "=", a]]
`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	var policy client.NotificationPolicy
	if err := UnmarshalInput(path, &policy); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(policy.GroupBy, []string{"alertname"}) {
		t.Errorf("GroupBy = %v, want [alertname]", policy.GroupBy)
	}
	if len(policy.Routes) != 1 || !reflect.DeepEqual(policy.Routes[0].ObjectMatchers, [][]string{{"team", "=", "a"}}) {
		t.Errorf("Routes = %+v, want one route with object_matchers [[team = a]]", policy.Routes)
	}
}

func TestUnmarshalInput_YAMLSilenceMatcherDefaultsToEqual(t *testing.T) {
	path := filepath.Join(t.TempDir(), "silence.yaml")
	data := `startsAt: 2026-11-01T00:00:00Z
matchers:
  - name: severity
    value: warning
    isRegex: false
`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	var req client.SilenceCreateRequest
	if err := UnmarshalInput(path, &req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.StartsAt != "2026-11-01T00:00:00Z" {
		t.Errorf("StartsAt = %q", req.StartsAt)
	}
	if len(req.Matchers) != 1 || !req.Matchers[0].IsEqual {
		t.Errorf("Matchers = %+v, want one matcher with IsEqual true", req.Matchers)
	}
}

func TestUnmarshalInput_YAMLNumbersInTextFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rule.yaml")
	data := `title: 2024
labels:
  priority: 1
`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	var rule client.AlertRule
	if err := UnmarshalInput(path, &rule); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Title != "2024" || rule.Labels["priority"] != "1" {
		t.Errorf("Title = %q, Labels = %v, want 2024 and priority=1", rule.Title, rule.Labels)
	}
}
