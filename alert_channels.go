package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AlertChannel represents a notification channel for alerts.
type AlertChannel struct {
	ID            string            `json:"id"`
	TeamID        string            `json:"team_id"`
	Type          string            `json:"type"`
	Name          string            `json:"name"`
	Config        map[string]string `json:"config"`
	Enabled       bool              `json:"enabled"`
	EmailVerified bool              `json:"email_verified"`
	Scope         string            `json:"scope"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

// CreateAlertChannelParams holds parameters for creating or updating an alert channel.
type CreateAlertChannelParams struct {
	Type    string            `json:"type"`
	Name    string            `json:"name"`
	Config  map[string]string `json:"config"`
	Enabled *bool             `json:"enabled,omitempty"`
}

// AlertEvent represents a single alert delivery event.
type AlertEvent struct {
	ID        string `json:"id"`
	MonitorID string `json:"monitor_id"`
	ChannelID string `json:"channel_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Delivered bool   `json:"delivered"`
	Error     string `json:"error"`
	CreatedAt string `json:"created_at"`
}

// MonitorChannelsResponse wraps the list of channel IDs linked to a monitor.
type MonitorChannelsResponse struct {
	ChannelIDs []string `json:"channel_ids"`
}

// alertChannelBasePath returns the base URL path for alert channel operations.
func alertChannelBasePath(teamID string) string {
	return "/api/v1/teams/" + teamID + "/alert-channels"
}

// alertChannelPath returns the URL path for a specific alert channel.
func alertChannelPath(teamID, channelID string) string {
	return alertChannelBasePath(teamID) + "/" + channelID
}

// monitorChannelsPath returns the URL path for a monitor's alert channel bindings.
func monitorChannelsPath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/alert-channels"
}

// CreateAlertChannel creates a new alert channel for a team.
func (c *Client) CreateAlertChannel(
	ctx context.Context, teamID string, params *CreateAlertChannelParams,
) (*AlertChannel, error) {
	body, err := c.do(ctx, http.MethodPost, alertChannelBasePath(teamID), params)
	if err != nil {
		return nil, err
	}
	var ch AlertChannel
	if err := json.Unmarshal(body, &ch); err != nil {
		return nil, fmt.Errorf("unmarshal alert channel: %w", err)
	}
	return &ch, nil
}

// ListAlertChannels returns all alert channels for a team.
func (c *Client) ListAlertChannels(ctx context.Context, teamID string) ([]AlertChannel, error) {
	body, err := c.do(ctx, http.MethodGet, alertChannelBasePath(teamID), nil)
	if err != nil {
		return nil, err
	}
	var channels []AlertChannel
	if err := json.Unmarshal(body, &channels); err != nil {
		return nil, fmt.Errorf("unmarshal alert channels: %w", err)
	}
	return channels, nil
}

// UpdateAlertChannel updates an existing alert channel.
func (c *Client) UpdateAlertChannel(
	ctx context.Context, teamID, channelID string, params *CreateAlertChannelParams,
) (*AlertChannel, error) {
	body, err := c.do(ctx, http.MethodPut, alertChannelPath(teamID, channelID), params)
	if err != nil {
		return nil, err
	}
	var ch AlertChannel
	if err := json.Unmarshal(body, &ch); err != nil {
		return nil, fmt.Errorf("unmarshal alert channel: %w", err)
	}
	return &ch, nil
}

// DeleteAlertChannel deletes an alert channel by ID.
func (c *Client) DeleteAlertChannel(ctx context.Context, teamID, channelID string) error {
	_, err := c.do(ctx, http.MethodDelete, alertChannelPath(teamID, channelID), nil)
	return err
}

// TestAlertChannel sends a test notification through the given channel.
func (c *Client) TestAlertChannel(ctx context.Context, teamID, channelID string) error {
	_, err := c.do(ctx, http.MethodPost, alertChannelPath(teamID, channelID)+"/test", nil)
	return err
}

// ListMonitorChannels returns the IDs of alert channels linked to a monitor.
func (c *Client) ListMonitorChannels(
	ctx context.Context, teamID, monitorID string,
) ([]string, error) {
	body, err := c.do(ctx, http.MethodGet, monitorChannelsPath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var resp MonitorChannelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal monitor channels: %w", err)
	}
	return resp.ChannelIDs, nil
}

// SetMonitorChannels replaces all alert channel bindings for a monitor.
func (c *Client) SetMonitorChannels(
	ctx context.Context, teamID, monitorID string, channelIDs []string,
) ([]string, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		monitorChannelsPath(teamID, monitorID),
		map[string][]string{"channel_ids": channelIDs},
	)
	if err != nil {
		return nil, err
	}
	var resp MonitorChannelsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal monitor channels: %w", err)
	}
	return resp.ChannelIDs, nil
}

// LinkMonitorChannel links a single alert channel to a monitor.
func (c *Client) LinkMonitorChannel(
	ctx context.Context, teamID, monitorID, channelID string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		monitorChannelsPath(teamID, monitorID)+"/"+channelID, nil,
	)
	return err
}

// UnlinkMonitorChannel removes a single alert channel from a monitor.
func (c *Client) UnlinkMonitorChannel(
	ctx context.Context, teamID, monitorID, channelID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		monitorChannelsPath(teamID, monitorID)+"/"+channelID, nil,
	)
	return err
}

// ListAlertEvents returns alert events for a monitor.
func (c *Client) ListAlertEvents(
	ctx context.Context, teamID, monitorID string,
) ([]AlertEvent, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		monitorPath(teamID, monitorID)+"/alert-events", nil,
	)
	if err != nil {
		return nil, err
	}
	var events []AlertEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("unmarshal alert events: %w", err)
	}
	return events, nil
}
