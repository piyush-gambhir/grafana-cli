package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type AlertRule struct {
	ID                          int64                  `json:"id"`
	UID                         string                 `json:"uid"`
	OrgID                       int64                  `json:"orgID"`
	FolderUID                   string                 `json:"folderUID"`
	RuleGroup                   string                 `json:"ruleGroup"`
	Title                       string                 `json:"title"`
	Condition                   string                 `json:"condition"`
	Data                        []AlertQuery           `json:"data"`
	NoDataState                 string                 `json:"noDataState"`
	ExecErrState                string                 `json:"execErrState"`
	For                         string                 `json:"for"`
	KeepFiringFor               string                 `json:"keep_firing_for,omitempty"`
	Annotations                 map[string]string      `json:"annotations,omitempty"`
	Labels                      map[string]string      `json:"labels,omitempty"`
	Updated                     string                 `json:"updated,omitempty"`
	Provenance                  string                 `json:"provenance,omitempty"`
	IsPaused                    bool                   `json:"isPaused"`
	NotificationSettings        map[string]interface{} `json:"notification_settings,omitempty"`
	Record                      map[string]interface{} `json:"record,omitempty"`
	MissingSeriesEvalsToResolve *int64                 `json:"missingSeriesEvalsToResolve,omitempty"`
}

type AlertQuery struct {
	RefID             string                 `json:"refId"`
	QueryType         string                 `json:"queryType,omitempty"`
	RelativeTimeRange RelativeTimeRange      `json:"relativeTimeRange"`
	DatasourceUID     string                 `json:"datasourceUid"`
	Model             map[string]interface{} `json:"model"`
}

type RelativeTimeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type AlertRuleGroup struct {
	Title     string      `json:"title"`
	FolderUID string      `json:"folderUid"`
	Interval  int64       `json:"interval"`
	Rules     []AlertRule `json:"rules"`
}

type ContactPoint struct {
	UID                   string                 `json:"uid"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	Settings              map[string]interface{} `json:"settings"`
	DisableResolveMessage bool                   `json:"disableResolveMessage"`
	Provenance            string                 `json:"provenance,omitempty"`
}

type NotificationPolicy struct {
	Receiver          string               `json:"receiver"`
	GroupBy           []string             `json:"group_by,omitempty"`
	GroupWait         string               `json:"group_wait,omitempty"`
	GroupInterval     string               `json:"group_interval,omitempty"`
	RepeatInterval    string               `json:"repeat_interval,omitempty"`
	ObjectMatchers    [][]string           `json:"object_matchers,omitempty"`
	Routes            []NotificationPolicy `json:"routes,omitempty"`
	Continue          bool                 `json:"continue,omitempty"`
	MuteTimeIntervals []string             `json:"mute_time_intervals,omitempty"`
}

type MuteTiming struct {
	Name          string         `json:"name"`
	TimeIntervals []TimeInterval `json:"time_intervals"`
	Provenance    string         `json:"provenance,omitempty"`
}

type TimeInterval struct {
	Times       []TimeRange `json:"times,omitempty"`
	Weekdays    []string    `json:"weekdays,omitempty"`
	DaysOfMonth []string    `json:"days_of_month,omitempty"`
	Months      []string    `json:"months,omitempty"`
	Years       []string    `json:"years,omitempty"`
	Location    string      `json:"location,omitempty"`
}

type TimeRange struct {
	StartMinute string `json:"start_time"`
	EndMinute   string `json:"end_time"`
}

type AlertTemplate struct {
	Name       string `json:"name"`
	Template   string `json:"template"`
	Provenance string `json:"provenance,omitempty"`
	Version    int64  `json:"version,omitempty"`
}

type AlertTemplatesResponse []AlertTemplate

type Silence struct {
	ID        string        `json:"id"`
	Status    SilenceStatus `json:"status"`
	Comment   string        `json:"comment"`
	CreatedBy string        `json:"createdBy"`
	StartsAt  string        `json:"startsAt"`
	EndsAt    string        `json:"endsAt"`
	Matchers  []Matcher     `json:"matchers"`
}

type SilenceStatus struct {
	State string `json:"state"`
}

type Matcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
	IsEqual bool   `json:"isEqual"`
}

func (m *Matcher) UnmarshalJSON(data []byte) error {
	type plain Matcher
	p := plain{IsEqual: true}
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	*m = Matcher(p)
	return nil
}

type SilenceCreateRequest struct {
	Comment   string    `json:"comment"`
	CreatedBy string    `json:"createdBy"`
	StartsAt  string    `json:"startsAt"`
	EndsAt    string    `json:"endsAt"`
	Matchers  []Matcher `json:"matchers"`
}

type SilenceCreateResponse struct {
	SilenceID string `json:"silenceID"`
}

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

func (c *Client) DeleteAlertRule(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/alert-rules/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

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

func (c *Client) GetContactPoint(ctx context.Context, uid string) (*ContactPoint, error) {
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

func (c *Client) UpdateContactPoint(ctx context.Context, uid string, cp ContactPoint) error {
	resp, err := c.Put(ctx, "/api/v1/provisioning/contact-points/"+uid, cp)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) DeleteContactPoint(ctx context.Context, uid string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/contact-points/"+uid)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

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

func (c *Client) UpdateNotificationPolicy(ctx context.Context, policy NotificationPolicy) error {
	resp, err := c.Put(ctx, "/api/v1/provisioning/policies", policy)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

func (c *Client) ResetNotificationPolicy(ctx context.Context) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/policies")
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

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

func (c *Client) DeleteMuteTiming(ctx context.Context, name string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/mute-timings/"+name)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

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

func (c *Client) DeleteAlertTemplate(ctx context.Context, name string) error {
	resp, err := c.Delete(ctx, "/api/v1/provisioning/templates/"+name)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}

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

func (c *Client) DeleteSilence(ctx context.Context, id string) error {
	resp, err := c.Delete(ctx, "/api/alertmanager/grafana/api/v2/silence/"+id)
	if err != nil {
		return err
	}
	return resp.JSON(nil)
}
