package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAlertRules(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/alert-rules" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]AlertRule{
			{UID: "rule1", Title: "High CPU"},
			{UID: "rule2", Title: "Low Memory"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	rules, err := c.ListAlertRules(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Title != "High CPU" {
		t.Errorf("got title %q, want %q", rules[0].Title, "High CPU")
	}
}

func TestGetAlertRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/alert-rules/rule1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(AlertRule{
			UID:       "rule1",
			Title:     "High CPU",
			Condition: "A",
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	rule, err := c.GetAlertRule(context.Background(), "rule1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Condition != "A" {
		t.Errorf("got condition %q, want %q", rule.Condition, "A")
	}
}

func TestCreateAlertRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/alert-rules" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req AlertRule
		json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "New Rule" {
			t.Errorf("got title %q, want %q", req.Title, "New Rule")
		}

		req.UID = "new-rule-uid"
		json.NewEncoder(w).Encode(req)
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateAlertRule(context.Background(), AlertRule{
		Title:     "New Rule",
		Condition: "A",
		FolderUID: "folder1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UID != "new-rule-uid" {
		t.Errorf("got UID %q, want %q", result.UID, "new-rule-uid")
	}
}

func TestDeleteAlertRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/alert-rules/rule1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.DeleteAlertRule(context.Background(), "rule1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListContactPoints(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/contact-points" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]ContactPoint{
			{UID: "cp1", Name: "Email", Type: "email"},
			{UID: "cp2", Name: "Slack", Type: "slack"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	points, err := c.ListContactPoints(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 contact points, got %d", len(points))
	}
	if points[1].Type != "slack" {
		t.Errorf("got type %q, want %q", points[1].Type, "slack")
	}
}

func TestCreateContactPoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/contact-points" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ContactPoint
		json.NewDecoder(r.Body).Decode(&req)

		req.UID = "new-cp"
		json.NewEncoder(w).Encode(req)
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	result, err := c.CreateContactPoint(context.Background(), ContactPoint{
		Name: "PagerDuty",
		Type: "pagerduty",
		Settings: map[string]interface{}{
			"integrationKey": "test-key",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UID != "new-cp" {
		t.Errorf("got UID %q, want %q", result.UID, "new-cp")
	}
}

func TestGetNotificationPolicy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/policies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(NotificationPolicy{
			Receiver: "grafana-default-email",
			GroupBy:  []string{"alertname"},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	policy, err := c.GetNotificationPolicy(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.Receiver != "grafana-default-email" {
		t.Errorf("got receiver %q, want %q", policy.Receiver, "grafana-default-email")
	}
	if len(policy.GroupBy) != 1 || policy.GroupBy[0] != "alertname" {
		t.Errorf("got group_by %v, want [alertname]", policy.GroupBy)
	}
}

func TestUpdateNotificationPolicy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/policies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req NotificationPolicy
		json.NewDecoder(r.Body).Decode(&req)
		if req.Receiver != "slack" {
			t.Errorf("got receiver %q, want %q", req.Receiver, "slack")
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Notification policies updated"})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	err := c.UpdateNotificationPolicy(context.Background(), NotificationPolicy{
		Receiver: "slack",
		GroupBy:  []string{"alertname", "cluster"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListMuteTimings(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/provisioning/mute-timings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]MuteTiming{
			{Name: "weekends", TimeIntervals: []TimeInterval{
				{Weekdays: []string{"saturday", "sunday"}},
			}},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	timings, err := c.ListMuteTimings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(timings) != 1 {
		t.Fatalf("expected 1 mute timing, got %d", len(timings))
	}
	if timings[0].Name != "weekends" {
		t.Errorf("got name %q, want %q", timings[0].Name, "weekends")
	}
}

func TestListSilences(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/alertmanager/grafana/api/v2/silences" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Silence{
			{
				ID:      "silence1",
				Comment: "Maintenance window",
				Status:  SilenceStatus{State: "active"},
				Matchers: []Matcher{
					{Name: "alertname", Value: "HighCPU", IsEqual: true},
				},
			},
		})
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}
	silences, err := c.ListSilences(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(silences) != 1 {
		t.Fatalf("expected 1 silence, got %d", len(silences))
	}
	if silences[0].Status.State != "active" {
		t.Errorf("got state %q, want %q", silences[0].Status.State, "active")
	}
	if silences[0].Comment != "Maintenance window" {
		t.Errorf("got comment %q, want %q", silences[0].Comment, "Maintenance window")
	}
}
