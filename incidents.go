package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Incident represents an account-level incident.
type AccountIncident struct {
	ID                 string   `json:"id"`
	AccountID          string   `json:"account_id"`
	Name               string   `json:"name"`
	Status             string   `json:"status"`
	Severity           string   `json:"severity"`
	Summary            string   `json:"summary"`
	CreatedBy          string   `json:"created_by"`
	OwnerID            string   `json:"owner_id,omitempty"`
	DeclaredAt         string   `json:"declared_at"`
	ResolvedAt         string   `json:"resolved_at,omitempty"`
	DurationSeconds    *int     `json:"duration_seconds,omitempty"`
	IsPrivate          bool     `json:"is_private"`
	Tags               []string `json:"tags"`
	Source             string   `json:"source"`
	MonitorIDs         []string `json:"monitor_ids,omitempty"`
	StatusPageIDs      []string `json:"status_page_ids,omitempty"`
	ServiceIDs         []string `json:"service_ids,omitempty"`
	EscalationPolicyID string   `json:"escalation_policy_id,omitempty"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

// IncidentUpdate represents a timeline update on an incident.
type AccountIncidentUpdate struct {
	ID                string `json:"id"`
	IncidentID        string `json:"incident_id"`
	Status            string `json:"status"`
	Message           string `json:"message"`
	CreatedBy         string `json:"created_by"`
	NotifySubscribers bool   `json:"notify_subscribers"`
	CreatedAt         string `json:"created_at"`
}

// EscalationState represents the current escalation state of an incident.
type EscalationState struct {
	Status         string `json:"status"`
	CurrentLevel   int    `json:"current_level"`
	AcknowledgedBy string `json:"acknowledged_by,omitempty"`
	AcknowledgedAt string `json:"acknowledged_at,omitempty"`
}

// EscalationEvent represents an event in the escalation timeline.
type EscalationEvent struct {
	Level      int    `json:"level"`
	UserID     string `json:"user_id"`
	ScheduleID string `json:"schedule_id"`
	Trigger    string `json:"trigger"`
	CreatedAt  string `json:"created_at"`
}

// IncidentWithDetails is an incident with its updates, escalation state, and events.
type AccountIncidentWithDetails struct {
	AccountIncident
	Updates          []AccountIncidentUpdate `json:"updates"`
	EscalationState  *EscalationState        `json:"escalation_state,omitempty"`
	EscalationEvents []EscalationEvent       `json:"escalation_events,omitempty"`
}

// CreateIncidentParams holds parameters for creating an incident.
type CreateAccountIncidentParams struct {
	Name                      string   `json:"name"`
	Severity                  string   `json:"severity,omitempty"`
	Summary                   string   `json:"summary,omitempty"`
	IsPrivate                 bool     `json:"is_private,omitempty"`
	Tags                      []string `json:"tags,omitempty"`
	MonitorIDs                []string `json:"monitor_ids,omitempty"`
	StatusPageIDs             []string `json:"status_page_ids,omitempty"`
	ServiceIDs                []string `json:"service_ids,omitempty"`
	PrimaryEscalationPolicyID string   `json:"primary_escalation_policy_id,omitempty"`
	NoEscalation              bool     `json:"no_escalation,omitempty"`
}

// UpdateIncidentParams holds parameters for updating an incident.
type UpdateAccountIncidentParams struct {
	Name      *string   `json:"name,omitempty"`
	Status    *string   `json:"status,omitempty"`
	Severity  *string   `json:"severity,omitempty"`
	Summary   *string   `json:"summary,omitempty"`
	OwnerID   *string   `json:"owner_id,omitempty"`
	IsPrivate *bool     `json:"is_private,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// PostIncidentUpdateParams holds parameters for posting an incident timeline update.
type PostAccountIncidentUpdateParams struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	NotifySubscribers bool   `json:"notify_subscribers"`
}

// ListIncidentsParams holds optional query parameters for listing incidents.
type ListIncidentsParams struct {
	Status    string
	Severity  string
	Tag       string
	ServiceID string
	From      string // RFC3339
	To        string // RFC3339
	Limit     int
	Offset    int
}

// IncidentAnalytics holds incident response metrics for a time range.
type AccountIncidentAnalytics struct {
	MTTR            *float64            `json:"mttr_seconds"`
	MTTF            *float64            `json:"mttf_seconds"`
	IncidentCount   int                 `json:"incident_count"`
	BySeverity      map[string]int      `json:"by_severity"`
	OpenActionItems int                 `json:"open_action_items"`
	MTTRBySeverity  map[string]*float64 `json:"mttr_by_severity"`
	UptimePct       *float64            `json:"uptime_pct,omitempty"`
}

// incidentBasePath returns the base URL path for incident operations.
func incidentBasePath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/incidents"
}

// incidentPath returns the URL path for a specific incident.
func incidentPath(accountID, incidentID string) string {
	return incidentBasePath(accountID) + "/" + incidentID
}

// CreateIncident creates a new incident for an account.
func (c *Client) CreateAccountIncident(
	ctx context.Context, accountID string, params *CreateAccountIncidentParams,
) (*AccountIncident, error) {
	body, err := c.do(ctx, http.MethodPost, incidentBasePath(accountID), params)
	if err != nil {
		return nil, err
	}
	var inc AccountIncident
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal account incident: %w", err)
	}
	return &inc, nil
}

// GetIncident returns a single incident with details (updates, escalation state/events).
func (c *Client) GetAccountIncident(
	ctx context.Context, accountID, incidentID string,
) (*AccountIncidentWithDetails, error) {
	body, err := c.do(ctx, http.MethodGet, incidentPath(accountID, incidentID), nil)
	if err != nil {
		return nil, err
	}
	var inc AccountIncidentWithDetails
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal account incident: %w", err)
	}
	return &inc, nil
}

// ListIncidents returns incidents for an account with optional filters.
func (c *Client) ListAccountIncidents(
	ctx context.Context, accountID string, params *ListIncidentsParams,
) ([]AccountIncident, error) {
	path := incidentBasePath(accountID)
	if params != nil {
		q := url.Values{}
		if params.Status != "" {
			q.Set("status", params.Status)
		}
		if params.Severity != "" {
			q.Set("severity", params.Severity)
		}
		if params.Tag != "" {
			q.Set("tag", params.Tag)
		}
		if params.ServiceID != "" {
			q.Set("service_id", params.ServiceID)
		}
		if params.From != "" {
			q.Set("from", params.From)
		}
		if params.To != "" {
			q.Set("to", params.To)
		}
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", fmt.Sprintf("%d", params.Offset))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var incidents []AccountIncident
	if err := json.Unmarshal(body, &incidents); err != nil {
		return nil, fmt.Errorf("unmarshal account incidents: %w", err)
	}
	return incidents, nil
}

// UpdateIncident updates an existing incident.
func (c *Client) UpdateAccountIncident(
	ctx context.Context, accountID, incidentID string, params *UpdateAccountIncidentParams,
) (*AccountIncident, error) {
	body, err := c.do(ctx, http.MethodPut, incidentPath(accountID, incidentID), params)
	if err != nil {
		return nil, err
	}
	var inc AccountIncident
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal account incident: %w", err)
	}
	return &inc, nil
}

// DeleteIncident deletes an incident by ID.
func (c *Client) DeleteAccountIncident(ctx context.Context, accountID, incidentID string) error {
	_, err := c.do(ctx, http.MethodDelete, incidentPath(accountID, incidentID), nil)
	return err
}

// PostIncidentUpdate posts a timeline update to an incident.
func (c *Client) PostAccountIncidentUpdate(
	ctx context.Context, accountID, incidentID string, params *PostAccountIncidentUpdateParams,
) (*AccountIncidentUpdate, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		incidentPath(accountID, incidentID)+"/updates", params,
	)
	if err != nil {
		return nil, err
	}
	var u AccountIncidentUpdate
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal account incident update: %w", err)
	}
	return &u, nil
}

// AcknowledgeIncident acknowledges an incident escalation.
func (c *Client) AcknowledgeAccountIncident(
	ctx context.Context, accountID, incidentID string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		incidentPath(accountID, incidentID)+"/acknowledge", nil,
	)
	return err
}

// LinkIncidentMonitors links monitors to an incident.
func (c *Client) LinkAccountIncidentMonitors(
	ctx context.Context, accountID, incidentID string, monitorIDs []string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		incidentPath(accountID, incidentID)+"/monitors",
		map[string][]string{"monitor_ids": monitorIDs},
	)
	return err
}

// UnlinkIncidentMonitor removes a monitor from an incident.
func (c *Client) UnlinkAccountIncidentMonitor(
	ctx context.Context, accountID, incidentID, monitorID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		incidentPath(accountID, incidentID)+"/monitors/"+monitorID, nil,
	)
	return err
}

// LinkIncidentStatusPages publishes an incident to status pages.
func (c *Client) LinkAccountIncidentStatusPages(
	ctx context.Context, accountID, incidentID string, statusPageIDs []string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		incidentPath(accountID, incidentID)+"/status-pages",
		map[string][]string{"status_page_ids": statusPageIDs},
	)
	return err
}

// UnlinkIncidentStatusPage removes an incident from a status page.
func (c *Client) UnlinkAccountIncidentStatusPage(
	ctx context.Context, accountID, incidentID, pageID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		incidentPath(accountID, incidentID)+"/status-pages/"+pageID, nil,
	)
	return err
}

// GetIncidentAnalytics returns incident response metrics for a time range.
func (c *Client) GetAccountIncidentAnalytics(
	ctx context.Context, accountID string, from, to string, serviceID string,
) (*AccountIncidentAnalytics, error) {
	q := url.Values{}
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	if serviceID != "" {
		q.Set("service_id", serviceID)
	}
	path := incidentBasePath(accountID) + "/analytics"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var analytics AccountIncidentAnalytics
	if err := json.Unmarshal(body, &analytics); err != nil {
		return nil, fmt.Errorf("unmarshal account incident analytics: %w", err)
	}
	return &analytics, nil
}
