package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ConsoleMessage represents a captured browser console message.
type ConsoleMessage struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	URL    string `json:"url"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// StepResult represents the result of a single browser step.
type StepResult struct {
	Action        string `json:"action"`
	Name          string `json:"name,omitempty"`
	Status        string `json:"status"`
	DurationMs    int64  `json:"duration_ms"`
	Error         string `json:"error,omitempty"`
	ScreenshotURL string `json:"screenshot_url,omitempty"`
}

// BrowserProbe represents a single browser probe result.
type BrowserProbe struct {
	ID                  string           `json:"id"`
	CheckerID           string           `json:"checker_id,omitempty"`
	CheckerPublicIP     string           `json:"checker_public_ip,omitempty"`
	Status              int              `json:"status"`
	Error               string           `json:"error"`
	TotalMs             int64            `json:"total_ms"`
	DomContentLoadedMs  int64            `json:"dom_content_loaded_ms"`
	LoadEventMs         int64            `json:"load_event_ms"`
	DomInteractiveMs    int64            `json:"dom_interactive_ms"`
	LcpMs               float64          `json:"lcp_ms"`
	FcpMs               float64          `json:"fcp_ms"`
	Cls                 float64          `json:"cls"`
	TtfbMs              float64          `json:"ttfb_ms"`
	ResourceCount       int              `json:"resource_count"`
	ResourceErrorCount  int              `json:"resource_error_count"`
	ResourceTotalBytes  int64            `json:"resource_total_bytes"`
	ResourceStatus1xx   int              `json:"resource_status_1xx,omitempty"`
	ResourceStatus2xx   int              `json:"resource_status_2xx,omitempty"`
	ResourceStatus3xx   int              `json:"resource_status_3xx,omitempty"`
	ResourceStatus4xx   int              `json:"resource_status_4xx,omitempty"`
	ResourceStatus5xx   int              `json:"resource_status_5xx,omitempty"`
	HarURL              string           `json:"har_url,omitempty"`
	ConsoleErrorCount   int              `json:"console_error_count"`
	ConsoleWarningCount int              `json:"console_warning_count"`
	ConsoleMessages     []ConsoleMessage `json:"console_messages,omitempty"`
	ScreenshotURL       string           `json:"screenshot_url,omitempty"`
	UserAgent           string           `json:"user_agent,omitempty"`
	StepResults         []StepResult     `json:"step_results,omitempty"`
	CheckedAt           string           `json:"checked_at"`
}

// BrowserProbeList holds a list of browser probes.
type BrowserProbeList struct {
	Items []BrowserProbe `json:"items"`
}

// BrowserProbeListOptions configures browser probe listing.
type BrowserProbeListOptions struct {
	From  string // RFC3339
	To    string // RFC3339
	Limit int
}

// BrowserProbeAggBucket represents a single time bucket in browser probe aggregation.
type BrowserProbeAggBucket struct {
	Timestamp          string  `json:"timestamp"`
	ProbeCount         int     `json:"probe_count"`
	ErrorCount         int     `json:"error_count"`
	TotalMs            float64 `json:"total_ms"`
	LcpMs              float64 `json:"lcp_ms"`
	FcpMs              float64 `json:"fcp_ms"`
	Cls                float64 `json:"cls"`
	TtfbMs             float64 `json:"ttfb_ms"`
	ResourceCount      float64 `json:"resource_count"`
	ResourceErrorCount float64 `json:"resource_error_count"`
}

// BrowserProbeAggregation holds aggregated browser probe data.
type BrowserProbeAggregation struct {
	Buckets []BrowserProbeAggBucket `json:"buckets"`
}

func browserProbeBasePath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/browser-probes"
}

// ListBrowserProbes returns browser probes for a monitor.
func (c *Client) ListBrowserProbes(
	ctx context.Context, teamID, monitorID string, opts BrowserProbeListOptions,
) (*BrowserProbeList, error) {
	path := browserProbeBasePath(teamID, monitorID)
	sep := "?"
	if opts.From != "" {
		path += sep + "from=" + opts.From
		sep = "&"
	}
	if opts.To != "" {
		path += sep + "to=" + opts.To
		sep = "&"
	}
	if opts.Limit > 0 {
		path += sep + "limit=" + strconv.Itoa(opts.Limit)
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var list BrowserProbeList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("unmarshal browser probe list: %w", err)
	}
	return &list, nil
}

// GetBrowserProbe returns a single browser probe with full details.
func (c *Client) GetBrowserProbe(
	ctx context.Context, teamID, monitorID, probeID string,
) (*BrowserProbe, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		browserProbeBasePath(teamID, monitorID)+"/"+probeID, nil,
	)
	if err != nil {
		return nil, err
	}
	var bp BrowserProbe
	if err := json.Unmarshal(body, &bp); err != nil {
		return nil, fmt.Errorf("unmarshal browser probe: %w", err)
	}
	return &bp, nil
}

// AggregateBrowserProbes returns aggregated browser probe metrics.
func (c *Client) AggregateBrowserProbes(
	ctx context.Context, teamID, monitorID string,
	from, to, step string,
) (*BrowserProbeAggregation, error) {
	path := browserProbeBasePath(teamID, monitorID) + "/aggregate" +
		"?from=" + from + "&to=" + to
	if step != "" {
		path += "&step=" + step
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var agg BrowserProbeAggregation
	if err := json.Unmarshal(body, &agg); err != nil {
		return nil, fmt.Errorf("unmarshal browser probe aggregation: %w", err)
	}
	return &agg, nil
}

// DownloadBrowserScreenshot downloads the screenshot for a browser probe.
func (c *Client) DownloadBrowserScreenshot(
	ctx context.Context, teamID, monitorID, probeID string,
) ([]byte, error) {
	return c.do(
		ctx, http.MethodGet,
		browserProbeBasePath(teamID, monitorID)+"/"+probeID+"/screenshot", nil,
	)
}

// DownloadBrowserHAR downloads the HAR file for a browser probe.
func (c *Client) DownloadBrowserHAR(
	ctx context.Context, teamID, monitorID, probeID string,
) ([]byte, error) {
	return c.do(
		ctx, http.MethodGet,
		browserProbeBasePath(teamID, monitorID)+"/"+probeID+"/har", nil,
	)
}
