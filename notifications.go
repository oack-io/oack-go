package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// NotificationDefaults holds the default notification channel IDs for an account.
type NotificationDefaults struct {
	ChannelIDs []string `json:"channel_ids"`
}

// MonitorNotification holds the notification channel IDs for a specific monitor.
type MonitorNotification struct {
	MonitorID  string   `json:"monitor_id"`
	ChannelIDs []string `json:"channel_ids"`
}

// GetNotificationDefaults returns the default notification channels for an account.
func (c *Client) GetNotificationDefaults(
	ctx context.Context, accountID string,
) (*NotificationDefaults, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/me/accounts/"+accountID+"/notification-defaults", nil,
	)
	if err != nil {
		return nil, err
	}
	var nd NotificationDefaults
	if err := json.Unmarshal(body, &nd); err != nil {
		return nil, fmt.Errorf("unmarshal notification defaults: %w", err)
	}
	return &nd, nil
}

// SetNotificationDefaults replaces the default notification channels for an account.
func (c *Client) SetNotificationDefaults(
	ctx context.Context, accountID string, channelIDs []string,
) (*NotificationDefaults, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		"/api/v1/me/accounts/"+accountID+"/notification-defaults",
		map[string][]string{"channel_ids": channelIDs},
	)
	if err != nil {
		return nil, err
	}
	var nd NotificationDefaults
	if err := json.Unmarshal(body, &nd); err != nil {
		return nil, fmt.Errorf("unmarshal notification defaults: %w", err)
	}
	return &nd, nil
}

// CopyAlertChannels copies alert channels from one account to another.
func (c *Client) CopyAlertChannels(
	ctx context.Context, fromAccountID, toAccountID string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		"/api/v1/me/alert-channels/copy",
		map[string]string{
			"from_account_id": fromAccountID,
			"to_account_id":   toAccountID,
		},
	)
	return err
}

// monitorNotificationPath returns the URL path for a monitor's notification settings.
func monitorNotificationPath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/my/notifications"
}

// GetMonitorNotification returns the notification settings for a monitor.
func (c *Client) GetMonitorNotification(
	ctx context.Context, teamID, monitorID string,
) (*MonitorNotification, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		monitorNotificationPath(teamID, monitorID), nil,
	)
	if err != nil {
		return nil, err
	}
	var mn MonitorNotification
	if err := json.Unmarshal(body, &mn); err != nil {
		return nil, fmt.Errorf("unmarshal monitor notification: %w", err)
	}
	return &mn, nil
}

// SetMonitorNotification replaces the notification channels for a monitor.
func (c *Client) SetMonitorNotification(
	ctx context.Context, teamID, monitorID string, channelIDs []string,
) (*MonitorNotification, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		monitorNotificationPath(teamID, monitorID),
		map[string][]string{"channel_ids": channelIDs},
	)
	if err != nil {
		return nil, err
	}
	var mn MonitorNotification
	if err := json.Unmarshal(body, &mn); err != nil {
		return nil, fmt.Errorf("unmarshal monitor notification: %w", err)
	}
	return &mn, nil
}

// RemoveMonitorNotification removes the notification settings for a monitor.
func (c *Client) RemoveMonitorNotification(
	ctx context.Context, teamID, monitorID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		monitorNotificationPath(teamID, monitorID), nil,
	)
	return err
}
