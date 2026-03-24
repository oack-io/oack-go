package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Monitor represents a monitored endpoint.
type Monitor struct {
	ID                      string            `json:"id"`
	TeamID                  string            `json:"team_id"`
	Name                    string            `json:"name"`
	URL                     string            `json:"url"`
	Status                  string            `json:"status"`
	TimeoutMs               int64             `json:"timeout_ms"`
	CheckIntervalMs         int64             `json:"check_interval_ms"`
	HTTPMethod              string            `json:"http_method"`
	HTTPVersion             string            `json:"http_version"`
	Headers                 map[string]string `json:"headers"`
	FollowRedirects         bool              `json:"follow_redirects"`
	AllowedStatusCodes      []string          `json:"allowed_status_codes"`
	FailureThreshold        int               `json:"failure_threshold"`
	LatencyThresholdMs      int               `json:"latency_threshold_ms"`
	SSLExpiryEnabled        bool              `json:"ssl_expiry_enabled"`
	SSLExpiryThresholds     []int             `json:"ssl_expiry_thresholds"`
	DomainExpiryEnabled     bool              `json:"domain_expiry_enabled"`
	DomainExpiryThresholds  []int             `json:"domain_expiry_thresholds"`
	UptimeThresholdGood     float64           `json:"uptime_threshold_good"`
	UptimeThresholdDegraded float64           `json:"uptime_threshold_degraded"`
	UptimeThresholdCritical float64           `json:"uptime_threshold_critical"`
	CheckerRegion           string            `json:"checker_region"`
	CheckerCountry          string            `json:"checker_country"`
	ResolveOverrideIP       string            `json:"resolve_override_ip"`
	HealthStatus            string            `json:"health_status"`
	HealthDownReason        string            `json:"health_down_reason"`
	ConsecutiveFailures     int               `json:"consecutive_failures"`
	ConsecutiveSuccesses    int               `json:"consecutive_successes"`
	HealthChangedAt         *string           `json:"health_changed_at"`
	IsDebugEnabled          bool              `json:"is_debug_enabled"`
	DebugExpiresAt          *string           `json:"debug_expires_at"`
	CheckerID               string            `json:"checker_id"`
	CFZoneIntegrationID     string            `json:"cf_zone_integration_id"`
	CreatedBy               string            `json:"created_by"`
	CreatedAt               string            `json:"created_at"`
	UpdatedAt               string            `json:"updated_at"`
}

// CreateMonitorParams holds parameters for creating or updating a monitor.
type CreateMonitorParams struct {
	Name                    string            `json:"name"`
	URL                     string            `json:"url"`
	CheckIntervalMs         int64             `json:"check_interval_ms,omitempty"`
	TimeoutMs               int64             `json:"timeout_ms,omitempty"`
	HTTPMethod              string            `json:"http_method,omitempty"`
	HTTPVersion             string            `json:"http_version,omitempty"`
	Headers                 map[string]string `json:"headers,omitempty"`
	FollowRedirects         *bool             `json:"follow_redirects,omitempty"`
	AllowedStatusCodes      []string          `json:"allowed_status_codes,omitempty"`
	FailureThreshold        int               `json:"failure_threshold,omitempty"`
	LatencyThresholdMs      int               `json:"latency_threshold_ms,omitempty"`
	SSLExpiryEnabled        *bool             `json:"ssl_expiry_enabled,omitempty"`
	SSLExpiryThresholds     []int             `json:"ssl_expiry_thresholds,omitempty"`
	DomainExpiryEnabled     *bool             `json:"domain_expiry_enabled,omitempty"`
	DomainExpiryThresholds  []int             `json:"domain_expiry_thresholds,omitempty"`
	UptimeThresholdGood     *float64          `json:"uptime_threshold_good,omitempty"`
	UptimeThresholdDegraded *float64          `json:"uptime_threshold_degraded,omitempty"`
	UptimeThresholdCritical *float64          `json:"uptime_threshold_critical,omitempty"`
	CheckerRegion           string            `json:"checker_region,omitempty"`
	CheckerCountry          string            `json:"checker_country,omitempty"`
	ResolveOverrideIP       string            `json:"resolve_override_ip,omitempty"`
	Status                  string            `json:"status,omitempty"`
}

// monitorBasePath returns the base URL path for monitor operations.
func monitorBasePath(teamID string) string {
	return "/api/v1/teams/" + teamID + "/monitors"
}

// monitorPath returns the URL path for a specific monitor.
func monitorPath(teamID, monitorID string) string {
	return monitorBasePath(teamID) + "/" + monitorID
}

// CreateMonitor creates a new monitor for a team.
func (c *Client) CreateMonitor(
	ctx context.Context, teamID string, params *CreateMonitorParams,
) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodPost, monitorBasePath(teamID), params)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// ListMonitors returns all monitors for a team.
func (c *Client) ListMonitors(ctx context.Context, teamID string) ([]Monitor, error) {
	body, err := c.do(ctx, http.MethodGet, monitorBasePath(teamID), nil)
	if err != nil {
		return nil, err
	}
	var monitors []Monitor
	if err := json.Unmarshal(body, &monitors); err != nil {
		return nil, fmt.Errorf("unmarshal monitors: %w", err)
	}
	return monitors, nil
}

// GetMonitor returns a single monitor by ID.
func (c *Client) GetMonitor(ctx context.Context, teamID, monitorID string) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodGet, monitorPath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// UpdateMonitor updates an existing monitor.
func (c *Client) UpdateMonitor(
	ctx context.Context, teamID, monitorID string, params *CreateMonitorParams,
) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodPut, monitorPath(teamID, monitorID), params)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// DeleteMonitor deletes a monitor by ID.
func (c *Client) DeleteMonitor(ctx context.Context, teamID, monitorID string) error {
	_, err := c.do(ctx, http.MethodDelete, monitorPath(teamID, monitorID), nil)
	return err
}

// PauseMonitor pauses a monitor's checks.
func (c *Client) PauseMonitor(ctx context.Context, teamID, monitorID string) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodPost, monitorPath(teamID, monitorID)+"/pause", nil)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// UnpauseMonitor resumes a paused monitor's checks.
func (c *Client) UnpauseMonitor(ctx context.Context, teamID, monitorID string) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodPost, monitorPath(teamID, monitorID)+"/unpause", nil)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// DuplicateMonitor creates a copy of an existing monitor.
func (c *Client) DuplicateMonitor(
	ctx context.Context, teamID, monitorID string,
) (*Monitor, error) {
	body, err := c.do(ctx, http.MethodPost, monitorPath(teamID, monitorID)+"/duplicate", nil)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// MoveMonitor moves a monitor to a different team.
func (c *Client) MoveMonitor(
	ctx context.Context, teamID, monitorID, targetTeamID string,
) (*Monitor, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		monitorPath(teamID, monitorID)+"/move",
		map[string]string{"target_team_id": targetTeamID},
	)
	if err != nil {
		return nil, err
	}
	var m Monitor
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor: %w", err)
	}
	return &m, nil
}

// TestMonitorAlert triggers a test alert for a monitor.
func (c *Client) TestMonitorAlert(ctx context.Context, teamID, monitorID string) error {
	_, err := c.do(ctx, http.MethodPost, monitorPath(teamID, monitorID)+"/test-alert", nil)
	return err
}
