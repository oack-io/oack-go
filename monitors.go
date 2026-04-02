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

	// Multi-location fields.
	AggregateFailureMode  string            `json:"aggregate_failure_mode,omitempty"`
	AggregateFailureCount int               `json:"aggregate_failure_count,omitempty"`
	Locations             []MonitorLocation `json:"locations,omitempty"`

	// Monitor type and browser config.
	Type          string         `json:"type"`
	BrowserConfig *BrowserConfig `json:"browser_config,omitempty"`
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

	// Multi-location fields.
	Locations             []LocationParams `json:"locations,omitempty"`
	AggregateFailureMode  string           `json:"aggregate_failure_mode,omitempty"`
	AggregateFailureCount int              `json:"aggregate_failure_count,omitempty"`

	// Monitor type and browser config.
	Type          string         `json:"type,omitempty"`
	BrowserConfig *BrowserConfig `json:"browser_config,omitempty"`
}

// BrowserConfig holds browser-specific monitor settings.
type BrowserConfig struct {
	ScreenshotEnabled      bool           `json:"screenshot_enabled"`
	ScreenshotFullPage     bool           `json:"screenshot_full_page"`
	ConsoleErrorThreshold  int            `json:"console_error_threshold"`
	ResourceErrorThreshold int            `json:"resource_error_threshold"`
	UserAgent              string         `json:"user_agent"`
	ViewportWidth          int            `json:"viewport_width"`
	ViewportHeight         int            `json:"viewport_height"`
	WaitUntil              string         `json:"wait_until"`
	ExtraWaitMs            int            `json:"extra_wait_ms"`
	Mode                   string         `json:"mode,omitempty"`
	Steps                  []BrowserStep  `json:"steps,omitempty"`
	Script                 string         `json:"script,omitempty"`
	ScriptEnv              []ScriptEnvVar `json:"script_env,omitempty"`
	SuiteURL               string         `json:"suite_url,omitempty"`
	DepsURL                string         `json:"deps_url,omitempty"`
	DepsHash               string         `json:"deps_hash,omitempty"`
	PwProject              string         `json:"pw_project,omitempty"`
	PwGrep                 string         `json:"pw_grep,omitempty"`
	PwTag                  string         `json:"pw_tag,omitempty"`
	SuiteGitSHA            string         `json:"suite_git_sha,omitempty"`
	SuiteGitBranch         string         `json:"suite_git_branch,omitempty"`
	SuiteGitOrigin         string         `json:"suite_git_origin,omitempty"`
	SuiteDeployHost        string         `json:"suite_deploy_host,omitempty"`
	SuiteUploadedAt        string         `json:"suite_uploaded_at,omitempty"`
	SuiteDeployedByID      string         `json:"suite_deployed_by_id,omitempty"`
	SuiteDeployedBy        string         `json:"suite_deployed_by,omitempty"`
	SuiteDeployedByImg     string         `json:"suite_deployed_by_img,omitempty"`
	SuiteDeployCmd         string         `json:"suite_deploy_cmd,omitempty"`
}

// BrowserStep represents a single step in a browser monitor step sequence.
type BrowserStep struct {
	Action       string `json:"action"`
	Selector     string `json:"selector,omitempty"`
	Value        string `json:"value,omitempty"`
	URL          string `json:"url,omitempty"`
	Attribute    string `json:"attribute,omitempty"`
	VariableName string `json:"variable_name,omitempty"`
	Name         string `json:"name,omitempty"`
	TimeoutMs    int    `json:"timeout_ms,omitempty"`
	WaitMs       int    `json:"wait_ms,omitempty"`
}

// ScriptEnvVar is an environment variable for script-mode browser monitors.
type ScriptEnvVar struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

// MonitorLocation represents a checker location assigned to a monitor.
type MonitorLocation struct {
	ID                string  `json:"id"`
	Label             string  `json:"label"`
	CheckerRegion     string  `json:"checker_region,omitempty"`
	CheckerID         *string `json:"checker_id,omitempty"`
	AssignedCheckerID *string `json:"assigned_checker_id,omitempty"`
	HealthStatus      string  `json:"health_status"`
	HealthDownReason  string  `json:"health_down_reason"`
	HealthChangedAt   *string `json:"health_changed_at,omitempty"`
}

// LocationParams holds parameters for specifying a monitor location.
type LocationParams struct {
	CheckerID     *string `json:"checker_id,omitempty"`
	CheckerRegion string  `json:"checker_region,omitempty"`
	Label         string  `json:"label,omitempty"`
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
