package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Trace represents a network trace for a monitor.
type Trace struct {
	ID          string  `json:"id"`
	MonitorID   string  `json:"monitor_id"`
	Status      string  `json:"status"`
	Result      any     `json:"result"`
	RequestedBy string  `json:"requested_by"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt *string `json:"completed_at"`
}

// traceBasePath returns the base URL path for trace operations.
func traceBasePath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/traces"
}

// ListTraces returns all traces for a monitor.
func (c *Client) ListTraces(
	ctx context.Context, teamID, monitorID string,
) ([]Trace, error) {
	body, err := c.do(ctx, http.MethodGet, traceBasePath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var traces []Trace
	if err := json.Unmarshal(body, &traces); err != nil {
		return nil, fmt.Errorf("unmarshal traces: %w", err)
	}
	return traces, nil
}

// RequestTrace initiates a new network trace for a monitor.
func (c *Client) RequestTrace(
	ctx context.Context, teamID, monitorID string,
) (*Trace, error) {
	body, err := c.do(ctx, http.MethodPost, traceBasePath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var t Trace
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal trace: %w", err)
	}
	return &t, nil
}
