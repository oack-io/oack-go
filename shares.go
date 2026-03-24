package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Share represents a public share link for a monitor.
type Share struct {
	ID        string  `json:"id"`
	MonitorID string  `json:"monitor_id"`
	Token     string  `json:"token"`
	ShareURL  string  `json:"share_url"`
	ExpiresAt *string `json:"expires_at"`
	CreatedBy string  `json:"created_by"`
	CreatedAt string  `json:"created_at"`
}

// CreateShareParams holds parameters for creating a share link.
type CreateShareParams struct {
	ExpiresAt *string `json:"expires_at,omitempty"`
}

// shareBasePath returns the base URL path for share operations.
func shareBasePath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/shares"
}

// CreateShare creates a new share link for a monitor.
func (c *Client) CreateShare(
	ctx context.Context, teamID, monitorID string, params CreateShareParams,
) (*Share, error) {
	body, err := c.do(ctx, http.MethodPost, shareBasePath(teamID, monitorID), params)
	if err != nil {
		return nil, err
	}
	var s Share
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal share: %w", err)
	}
	return &s, nil
}

// ListShares returns all share links for a monitor.
func (c *Client) ListShares(
	ctx context.Context, teamID, monitorID string,
) ([]Share, error) {
	body, err := c.do(ctx, http.MethodGet, shareBasePath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var shares []Share
	if err := json.Unmarshal(body, &shares); err != nil {
		return nil, fmt.Errorf("unmarshal shares: %w", err)
	}
	return shares, nil
}

// RevokeShare revokes a share link for a monitor.
func (c *Client) RevokeShare(ctx context.Context, teamID, monitorID, shareID string) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		shareBasePath(teamID, monitorID)+"/"+shareID, nil,
	)
	return err
}
