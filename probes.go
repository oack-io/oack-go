package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Probe represents a single probe result for a monitor.
type Probe struct {
	ID             string  `json:"id"`
	MonitorID      string  `json:"monitor_id"`
	CheckerID      string  `json:"checker_id"`
	CheckerRegion  string  `json:"checker_region"`
	CheckerCountry string  `json:"checker_country"`
	StatusCode     int     `json:"status_code"`
	ResponseTimeMs float64 `json:"response_time_ms"`
	DNSTimeMs      float64 `json:"dns_time_ms"`
	ConnectTimeMs  float64 `json:"connect_time_ms"`
	TLSTimeMs      float64 `json:"tls_time_ms"`
	TTFBMs         float64 `json:"ttfb_ms"`
	TransferTimeMs float64 `json:"transfer_time_ms"`
	Error          string  `json:"error"`
	IsUp           bool    `json:"is_up"`
	CreatedAt      string  `json:"created_at"`
}

// ProbeList is a paginated list of probes.
type ProbeList struct {
	Probes []Probe `json:"probes"`
	Total  int     `json:"total"`
}

// ProbeListOptions configures probe listing.
type ProbeListOptions struct {
	Limit  int
	Offset int
	IsUp   *bool
}

// ProbeAggBucket represents a single time bucket in a probe aggregation.
type ProbeAggBucket struct {
	Timestamp     string  `json:"timestamp"`
	AvgResponseMs float64 `json:"avg_response_ms"`
	MinResponseMs float64 `json:"min_response_ms"`
	MaxResponseMs float64 `json:"max_response_ms"`
	SuccessCount  int     `json:"success_count"`
	FailureCount  int     `json:"failure_count"`
	TotalCount    int     `json:"total_count"`
}

// ProbeAggregation holds aggregated probe data.
type ProbeAggregation struct {
	Buckets []ProbeAggBucket `json:"buckets"`
}

// probeBasePath returns the base URL path for probe operations.
func probeBasePath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/probes"
}

// ListProbes returns a paginated list of probes for a monitor.
func (c *Client) ListProbes(
	ctx context.Context, teamID, monitorID string, opts ProbeListOptions,
) (*ProbeList, error) {
	path := probeBasePath(teamID, monitorID)
	sep := "?"
	if opts.Limit > 0 {
		path += sep + "limit=" + strconv.Itoa(opts.Limit)
		sep = "&"
	}
	if opts.Offset > 0 {
		path += sep + "offset=" + strconv.Itoa(opts.Offset)
		sep = "&"
	}
	if opts.IsUp != nil {
		path += sep + "is_up=" + strconv.FormatBool(*opts.IsUp)
	}
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var list ProbeList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("unmarshal probe list: %w", err)
	}
	return &list, nil
}

// GetProbe returns a single probe by ID.
func (c *Client) GetProbe(
	ctx context.Context, teamID, monitorID, probeID string,
) (*Probe, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		probeBasePath(teamID, monitorID)+"/"+probeID, nil,
	)
	if err != nil {
		return nil, err
	}
	var p Probe
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("unmarshal probe: %w", err)
	}
	return &p, nil
}

// GetProbeDetails returns the raw JSON details for a probe.
func (c *Client) GetProbeDetails(
	ctx context.Context, teamID, monitorID, probeID string,
) (json.RawMessage, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		probeBasePath(teamID, monitorID)+"/"+probeID+"/details", nil,
	)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

// DownloadPcap downloads the pcap file for a probe.
func (c *Client) DownloadPcap(
	ctx context.Context, teamID, monitorID, probeID string,
) ([]byte, error) {
	return c.do(
		ctx, http.MethodGet,
		probeBasePath(teamID, monitorID)+"/"+probeID+"/pcap", nil,
	)
}

// AggregateProbes returns aggregated probe data for a monitor.
func (c *Client) AggregateProbes(
	ctx context.Context, teamID, monitorID string,
	from, to int64, step, agg string,
) (*ProbeAggregation, error) {
	path := probeBasePath(teamID, monitorID) + "/aggregate" +
		"?from=" + strconv.FormatInt(from, 10) +
		"&to=" + strconv.FormatInt(to, 10) +
		"&step=" + step +
		"&agg=" + agg
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var result ProbeAggregation
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal probe aggregation: %w", err)
	}
	return &result, nil
}
