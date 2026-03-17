package client

import (
	"context"
	"fmt"
)

// AlertRule represents a Grafana alerting rule.
type AlertRule struct {
	ID           int64                  `json:"id"`
	UID          string                 `json:"uid"`
	OrgID        int64                  `json:"orgID"`
	FolderUID    string                 `json:"folderUID"`
	RuleGroup    string                 `json:"ruleGroup"`
	Title        string                 `json:"title"`
	Condition    string                 `json:"condition"`
	Data         []AlertQuery           `json:"data"`
	NoDataState  string                 `json:"noDataState"`
	ExecErrState string                 `json:"execErrState"`
	For          string                 `json:"for"`
	Annotations  map[string]string      `json:"annotations,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	Updated      string                 `json:"updated,omitempty"`
	Provenance   string                 `json:"provenance,omitempty"`
}

// AlertQuery is a query within an alert rule.
type AlertQuery struct {
	RefID             string                 `json:"refId"`
	QueryType         string                 `json:"queryType,omitempty"`
	RelativeTimeRange RelativeTimeRange      `json:"relativeTimeRange"`
	DatasourceUID     string                 `json:"datasourceUid"`
	Model             map[string]interface{} `json:"model"`
}

// RelativeTimeRange represents a relative time range.
type RelativeTimeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// AlertRuleGroup represents a group of alert rules.
type AlertRuleGroup struct {
	Title     string      `json:"title"`
	FolderUID string      `json:"folderUid"`
	Interval  int64       `json:"interval"`
	Rules     []AlertRule `json:"rules"`
}

// ContactPoint represents a contact point in Grafana alerting.
type ContactPoint struct {
	UID                   string                 `json:"uid"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	Settings              map[string]interface{} `json:"settings"`
	DisableResolveMessage bool                   `json:"disableResolveMessage"`
	Provenance            string                 `json:"provenance,omitempty"`
}

// NotificationPolicy represents the notification policy tree.
type NotificationPolicy struct {
	Receiver       string                `json:"receiver"`
	GroupBy        []string              `json:"group_by,omitempty"`
	GroupWait      string                `json:"group_wait,omitempty"`
	GroupInterval  string                `json:"group_interval,omitempty"`
	RepeatInterval string                `json:"repeat_interval,omitempty"`
	ObjectMatchers [][]string            `json:"object_matchers,omitempty"`
	Routes         []NotificationPolicy  `json:"routes,omitempty"`
	Continue       bool                  `json:"continue,omitempty"`
	MuteTimeIntervals []string           `json:"mute_time_intervals,omitempty"`
}

// MuteTiming represents a mute timing.
type MuteTiming struct {
	Name          string         `json:"name"`
	TimeIntervals []TimeInterval `json:"time_intervals"`
	Provenance    string         `json:"provenance,omitempty"`
}

// TimeInterval represents a time interval within a mute timing.
type TimeInterval struct {
	Times       []TimeRange `json:"times,omitempty"`
	Weekdays    []string    `json:"weekdays,omitempty"`
	DaysOfMonth []string    `json:"days_of_month,omitempty"`
	Months      []string    `json:"months,omitempty"`
	Years       []string    `json:"years,omitempty"`
	Location    string      `json:"location,omitempty"`
}

// TimeRange represents a time range.
type TimeRange struct {
	StartMinute string `json:"start_time"`
	EndMinute   string `json:"end_time"`
}

// AlertTemplate represents a notification template.
type AlertTemplate struct {
	Name       string `json:"name"`
	Template   string `json:"template"`
	Provenance string `json:"provenance,omitempty"`
	Version    int64  `json:"version,omitempty"`
}

// AlertTemplatesResponse is the response for listing templates.
type AlertTemplatesResponse []AlertTemplate

// Silence represents an alerting silence.
type Silence struct {
	ID        string    `json:"id"`
	Status    SilenceStatus `json:"status"`
	Comment   string    `json:"comment"`
	CreatedBy string    `json:"createdBy"`
	StartsAt  string    `json:"startsAt"`
	EndsAt    string    `json:"endsAt"`
	Matchers  []Matcher `json:"matchers"`
}

// SilenceStatus represents the status of a silence.
type SilenceStatus struct {
	State string `json:"state"`
}

// Matcher represents a label matcher.
type Matcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
	IsEqual bool   `json:"isEqual"`
}

// SilenceCreateRequest is the request body for creating a silence.
type SilenceCreateRequest struct {
	Comment   string    `json:"comment"`
	CreatedBy string    `json:"createdBy"`
	StartsAt  string    `json:"startsAt"`
	EndsAt    string    `json:"endsAt"`
	Matchers  []Matcher `json:"matchers"`
}

// SilenceCreateResponse is the response from creating a silence.
type SilenceCreateResponse struct {
	SilenceID string `json:"silenceID"`
}

// ListAlertRules returns all alert rules.
func (c *Client) ListAlertRules(ctx context.Context) ([]AlertRule, error) {
	var results []AlertRule
	resp, err := c.Get(ctx, "/api/v1/provisioning/alert-rules")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetAlertRule gets an alert rule by UID.
func (c *Client) GetAlertRule(ctx context.Context, uid string) (*AlertRule, error) {
	var result AlertRule
	resp, err := c.Get(ctx, "/api/v1/provisioning/alert-rules/"+uid)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateAlertRule creates a new alert rule.
func (c *Client) CreateAlertRule(ctx context.Context, rule AlertRule) (*AlertRule, error) {
	var result AlertRule
	resp, err := c.Post(ctx, "/api/v1/provisioning/alert-rules", rule)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateAlertRule updates an existing alert rule.
func (c *Client) UpdateAlertRule(ctx context.Context, uid string, rule AlertRule) (*AlertRule, error) {
	var result AlertRule
	resp, err := c.Put(ctx, "/api/v1/provisioning/alert-rules/"+uid, rule)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAlertRule deletes an alert rule by UID.
func (c *Client) DeleteAlertRule(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/alert-rules/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ListContactPoints returns all contact points.
func (c *Client) ListContactPoints(ctx context.Context) ([]ContactPoint, error) {
	var results []ContactPoint
	resp, err := c.Get(ctx, "/api/v1/provisioning/contact-points")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetContactPoint returns a contact point by UID.
func (c *Client) GetContactPoint(ctx context.Context, uid string) (*ContactPoint, error) {
	// Grafana doesn't have a direct get-by-UID endpoint, so we list and filter.
	points, err := c.ListContactPoints(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range points {
		if p.UID == uid {
			return &p, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Message: fmt.Sprintf("contact point %q not found", uid)}
}

// CreateContactPoint creates a new contact point.
func (c *Client) CreateContactPoint(ctx context.Context, cp ContactPoint) (*ContactPoint, error) {
	var result ContactPoint
	resp, err := c.Post(ctx, "/api/v1/provisioning/contact-points", cp)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateContactPoint updates an existing contact point.
func (c *Client) UpdateContactPoint(ctx context.Context, uid string, cp ContactPoint) error {
	resp, err := c.Put(ctx, "/api/v1/provisioning/contact-points/"+uid, cp)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// DeleteContactPoint deletes a contact point by UID.
func (c *Client) DeleteContactPoint(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/contact-points/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// GetNotificationPolicy returns the notification policy tree.
func (c *Client) GetNotificationPolicy(ctx context.Context) (*NotificationPolicy, error) {
	var result NotificationPolicy
	resp, err := c.Get(ctx, "/api/v1/provisioning/policies")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateNotificationPolicy updates the notification policy tree.
func (c *Client) UpdateNotificationPolicy(ctx context.Context, policy NotificationPolicy) error {
	resp, err := c.Put(ctx, "/api/v1/provisioning/policies", policy)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ResetNotificationPolicy resets the notification policy to defaults.
func (c *Client) ResetNotificationPolicy(ctx context.Context) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/policies")
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ListMuteTimings returns all mute timings.
func (c *Client) ListMuteTimings(ctx context.Context) ([]MuteTiming, error) {
	var results []MuteTiming
	resp, err := c.Get(ctx, "/api/v1/provisioning/mute-timings")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetMuteTiming returns a mute timing by name.
func (c *Client) GetMuteTiming(ctx context.Context, name string) (*MuteTiming, error) {
	var result MuteTiming
	resp, err := c.Get(ctx, "/api/v1/provisioning/mute-timings/"+name)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateMuteTiming creates a new mute timing.
func (c *Client) CreateMuteTiming(ctx context.Context, mt MuteTiming) (*MuteTiming, error) {
	var result MuteTiming
	resp, err := c.Post(ctx, "/api/v1/provisioning/mute-timings", mt)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateMuteTiming updates an existing mute timing.
func (c *Client) UpdateMuteTiming(ctx context.Context, name string, mt MuteTiming) (*MuteTiming, error) {
	var result MuteTiming
	resp, err := c.Put(ctx, "/api/v1/provisioning/mute-timings/"+name, mt)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMuteTiming deletes a mute timing by name.
func (c *Client) DeleteMuteTiming(ctx context.Context, name string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/mute-timings/"+name)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ListAlertTemplates returns all notification templates.
func (c *Client) ListAlertTemplates(ctx context.Context) ([]AlertTemplate, error) {
	var results []AlertTemplate
	resp, err := c.Get(ctx, "/api/v1/provisioning/templates")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetAlertTemplate returns a notification template by name.
func (c *Client) GetAlertTemplate(ctx context.Context, name string) (*AlertTemplate, error) {
	var result AlertTemplate
	resp, err := c.Get(ctx, "/api/v1/provisioning/templates/"+name)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateAlertTemplate creates or updates a notification template.
func (c *Client) UpdateAlertTemplate(ctx context.Context, name string, tmpl AlertTemplate) (*AlertTemplate, error) {
	var result AlertTemplate
	resp, err := c.Put(ctx, "/api/v1/provisioning/templates/"+name, tmpl)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAlertTemplate deletes a notification template by name.
func (c *Client) DeleteAlertTemplate(ctx context.Context, name string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/templates/"+name)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

// ListSilences returns all silences.
func (c *Client) ListSilences(ctx context.Context) ([]Silence, error) {
	var results []Silence
	resp, err := c.Get(ctx, "/api/alertmanager/grafana/api/v2/silences")
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetSilence returns a silence by ID.
func (c *Client) GetSilence(ctx context.Context, id string) (*Silence, error) {
	var result Silence
	resp, err := c.Get(ctx, "/api/alertmanager/grafana/api/v2/silence/"+id)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSilence creates a new silence.
func (c *Client) CreateSilence(ctx context.Context, req SilenceCreateRequest) (*SilenceCreateResponse, error) {
	var result SilenceCreateResponse
	resp, err := c.Post(ctx, "/api/alertmanager/grafana/api/v2/silences", req)
	if err != nil {
		return nil, err
	}
	if err := resp.JSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteSilence deletes a silence by ID.
func (c *Client) DeleteSilence(ctx context.Context, id string) error {
	resp, err := c.Delete(ctx, "/api/alertmanager/grafana/api/v2/silence/"+id)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
