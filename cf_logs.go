package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// CFLogEntry represents a Cloudflare log entry for a probe.
type CFLogEntry struct {
	ID                   string `json:"id"`
	ProbeID              string `json:"probe_id"`
	CfRay                string `json:"cf_ray"`
	EdgeColoCode         string `json:"edge_colo_code"`
	CacheStatus          string `json:"cache_status"`
	EdgeResponseStatus   int    `json:"edge_response_status"`
	OriginResponseStatus int    `json:"origin_response_status"`
	CreatedAt            string `json:"created_at"`
}

// CFLogListOptions holds query parameters for listing CF logs.
type CFLogListOptions struct {
	From  *string
	To    *string
	Limit *int
}

func (o CFLogListOptions) queryString() string {
	v := url.Values{}
	if o.From != nil {
		v.Set("from", *o.From)
	}
	if o.To != nil {
		v.Set("to", *o.To)
	}
	if o.Limit != nil {
		v.Set("limit", strconv.Itoa(*o.Limit))
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// GetCFLog returns a single CF log entry for a probe.
func (c *Client) GetCFLog(
	ctx context.Context, teamID, monitorID, probeID string,
) (*CFLogEntry, error) {
	path := monitorPath(teamID, monitorID) + "/probes/" + probeID + "/cf-log"
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var entry CFLogEntry
	if err := json.Unmarshal(body, &entry); err != nil {
		return nil, fmt.Errorf("unmarshal cf log: %w", err)
	}
	return &entry, nil
}

// ListCFLogs returns CF logs for a monitor.
func (c *Client) ListCFLogs(
	ctx context.Context, teamID, monitorID string, opts CFLogListOptions,
) ([]CFLogEntry, error) {
	path := monitorPath(teamID, monitorID) + "/cf-logs" + opts.queryString()
	body, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var entries []CFLogEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal cf logs: %w", err)
	}
	return entries, nil
}
