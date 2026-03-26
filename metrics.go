package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// WindowMetrics holds aggregated metrics for a time window.
type WindowMetrics struct {
	Uptime        float64 `json:"uptime"`
	AvgResponseMs float64 `json:"avg_response_ms"`
	P95ResponseMs float64 `json:"p95_response_ms"`
	TotalProbes   int     `json:"total_probes"`
	SuccessProbes int     `json:"success_probes"`
	FailureProbes int     `json:"failure_probes"`
}

// MonitorMetrics holds metrics across multiple time windows.
type MonitorMetrics struct {
	Last24h WindowMetrics `json:"last_24h"`
	Last7d  WindowMetrics `json:"last_7d"`
	Last30d WindowMetrics `json:"last_30d"`
}

// ExpirationSSL holds SSL certificate expiration details.
type ExpirationSSL struct {
	ExpiresAt *string `json:"expires_at"`
	Issuer    string  `json:"issuer"`
	Subject   string  `json:"subject"`
	DaysLeft  *int    `json:"days_left"`
	Status    string  `json:"status"`
	CheckedAt *string `json:"checked_at"`
}

// ExpirationDomain holds domain registration expiration details.
type ExpirationDomain struct {
	ExpiresAt *string `json:"expires_at"`
	Registrar string  `json:"registrar"`
	DaysLeft  *int    `json:"days_left"`
	Status    string  `json:"status"`
	CheckedAt *string `json:"checked_at"`
}

// Expiration holds SSL and domain expiration information.
type Expiration struct {
	SSL    *ExpirationSSL    `json:"ssl"`
	Domain *ExpirationDomain `json:"domain"`
}

// TimelineEvent represents an event in a monitor's timeline.
type TimelineEvent struct {
	ID        string `json:"id"`
	MonitorID string `json:"monitor_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// ChartEvent represents an annotation event on a team's chart.
type ChartEvent struct {
	ID          string         `json:"id"`
	TeamID      string         `json:"team_id"`
	MonitorID   string         `json:"monitor_id"`
	Source      string         `json:"source"`
	Kind        string         `json:"kind"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Severity    string         `json:"severity"`
	ExternalID  string         `json:"external_id"`
	StartedAt   string         `json:"started_at"`
	EndedAt     *string        `json:"ended_at"`
	Metadata    map[string]any `json:"metadata"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// CreateChartEventParams holds parameters for creating a chart event.
type CreateChartEventParams struct {
	MonitorID   string         `json:"monitor_id,omitempty"`
	Kind        string         `json:"kind"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Severity    string         `json:"severity,omitempty"`
	StartedAt   string         `json:"started_at"`
	EndedAt     *string        `json:"ended_at,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// UpdateChartEventParams holds parameters for updating a chart event.
type UpdateChartEventParams struct {
	MonitorID   *string        `json:"monitor_id,omitempty"`
	Kind        *string        `json:"kind,omitempty"`
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Severity    *string        `json:"severity,omitempty"`
	StartedAt   *string        `json:"started_at,omitempty"`
	EndedAt     *string        `json:"ended_at,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// ChartEventListOptions configures chart event listing.
type ChartEventListOptions struct {
	From      string // Required, RFC3339
	To        string // Required, RFC3339
	MonitorID string
	Kind      string // Comma-separated
	Source    string // Comma-separated
}

// IngestChartEventParams holds parameters for the CI/CD event ingest endpoint.
type IngestChartEventParams struct {
	Kind      string         `json:"kind"`
	Title     string         `json:"title"`
	URL       string         `json:"url,omitempty"`
	StartedAt *string        `json:"started_at,omitempty"`
	EndedAt   *string        `json:"ended_at,omitempty"`
	MonitorID *string        `json:"monitor_id,omitempty"`
	Severity  string         `json:"severity,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// TimelineListOptions configures timeline listing.
type TimelineListOptions struct {
	Limit  int
	Offset int
}

// GetMonitorMetrics returns aggregated metrics for a monitor.
func (c *Client) GetMonitorMetrics(
	ctx context.Context, teamID, monitorID string,
) (*MonitorMetrics, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		monitorPath(teamID, monitorID)+"/metrics", nil,
	)
	if err != nil {
		return nil, err
	}
	var m MonitorMetrics
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal monitor metrics: %w", err)
	}
	return &m, nil
}

// GetMonitorExpiration returns SSL and domain expiration info for a monitor.
func (c *Client) GetMonitorExpiration(
	ctx context.Context, teamID, monitorID string,
) (*Expiration, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		monitorPath(teamID, monitorID)+"/expiration", nil,
	)
	if err != nil {
		return nil, err
	}
	var exp Expiration
	if err := json.Unmarshal(body, &exp); err != nil {
		return nil, fmt.Errorf("unmarshal expiration: %w", err)
	}
	return &exp, nil
}

// ListTimeline returns timeline events for a monitor.
func (c *Client) ListTimeline(
	ctx context.Context, teamID, monitorID string, opts TimelineListOptions,
) ([]TimelineEvent, error) {
	path := monitorPath(teamID, monitorID) + "/timeline"
	sep := "?"
	if opts.Limit > 0 {
		path += sep + "limit=" + strconv.Itoa(opts.Limit)
		sep = "&"
	}
	if opts.Offset > 0 {
		path += sep + "offset=" + strconv.Itoa(opts.Offset)
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var events []TimelineEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("unmarshal timeline events: %w", err)
	}
	return events, nil
}

// chartEventBasePath returns the base URL path for chart event operations.
func chartEventBasePath(teamID string) string {
	return "/api/v1/teams/" + teamID + "/events"
}

// CreateChartEvent creates a new chart event for a team.
func (c *Client) CreateChartEvent(
	ctx context.Context, teamID string, params CreateChartEventParams,
) (*ChartEvent, error) {
	body, err := c.do(ctx, http.MethodPost, chartEventBasePath(teamID), params)
	if err != nil {
		return nil, err
	}
	var ev ChartEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("unmarshal chart event: %w", err)
	}
	return &ev, nil
}

// ListChartEvents returns chart events for a team within a time range.
func (c *Client) ListChartEvents(
	ctx context.Context, teamID string, opts ChartEventListOptions,
) ([]ChartEvent, error) {
	path := chartEventBasePath(teamID) +
		"?from=" + opts.From + "&to=" + opts.To
	if opts.MonitorID != "" {
		path += "&monitor_id=" + opts.MonitorID
	}
	if opts.Kind != "" {
		path += "&kind=" + opts.Kind
	}
	if opts.Source != "" {
		path += "&source=" + opts.Source
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var events []ChartEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("unmarshal chart events: %w", err)
	}
	return events, nil
}

// UpdateChartEvent updates an existing chart event.
func (c *Client) UpdateChartEvent(
	ctx context.Context, teamID, eventID string, params UpdateChartEventParams,
) (*ChartEvent, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		chartEventBasePath(teamID)+"/"+eventID, params,
	)
	if err != nil {
		return nil, err
	}
	var ev ChartEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("unmarshal chart event: %w", err)
	}
	return &ev, nil
}

// DeleteChartEvent deletes a chart event by ID.
func (c *Client) DeleteChartEvent(ctx context.Context, teamID, eventID string) error {
	_, err := c.do(ctx, http.MethodDelete, chartEventBasePath(teamID)+"/"+eventID, nil)
	return err
}

// IngestChartEvent creates a chart event using API key authentication (CI/CD).
func (c *Client) IngestChartEvent(
	ctx context.Context, teamID string, params IngestChartEventParams,
) (*ChartEvent, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		chartEventBasePath(teamID)+"/ingest", params,
	)
	if err != nil {
		return nil, err
	}
	var ev ChartEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("unmarshal chart event: %w", err)
	}
	return &ev, nil
}
