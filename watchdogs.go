package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Watchdog links a monitor to a status page component for automatic incident management.
type Watchdog struct {
	ID                string `json:"id"`
	ComponentID       string `json:"component_id"`
	MonitorID         string `json:"monitor_id"`
	Severity          string `json:"severity"`
	AutoCreate        bool   `json:"auto_create"`
	AutoResolve       bool   `json:"auto_resolve"`
	NotifySubscribers bool   `json:"notify_subscribers"`
	TemplateID        string `json:"template_id"`
	CreatedAt         string `json:"created_at"`
}

// CreateWatchdogParams holds parameters for creating or updating a watchdog.
type CreateWatchdogParams struct {
	MonitorID         string `json:"monitor_id"`
	Severity          string `json:"severity"`
	AutoCreate        *bool  `json:"auto_create,omitempty"`
	AutoResolve       *bool  `json:"auto_resolve,omitempty"`
	NotifySubscribers *bool  `json:"notify_subscribers,omitempty"`
	TemplateID        string `json:"template_id,omitempty"`
}

func watchdogPath(accountID, pageID, compID string) string {
	return statusPagePath(accountID, pageID) + "/components/" + compID + "/watchdogs"
}

// CreateWatchdog creates a new watchdog on a status page component.
func (c *Client) CreateWatchdog(
	ctx context.Context, accountID, pageID, compID string, params *CreateWatchdogParams,
) (*Watchdog, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		watchdogPath(accountID, pageID, compID), params,
	)
	if err != nil {
		return nil, err
	}
	var w Watchdog
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("unmarshal watchdog: %w", err)
	}
	return &w, nil
}

// ListWatchdogs returns all watchdogs for a status page component.
func (c *Client) ListWatchdogs(
	ctx context.Context, accountID, pageID, compID string,
) ([]Watchdog, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		watchdogPath(accountID, pageID, compID), nil,
	)
	if err != nil {
		return nil, err
	}
	var watchdogs []Watchdog
	if err := json.Unmarshal(body, &watchdogs); err != nil {
		return nil, fmt.Errorf("unmarshal watchdogs: %w", err)
	}
	return watchdogs, nil
}

// UpdateWatchdog updates an existing watchdog.
func (c *Client) UpdateWatchdog(
	ctx context.Context,
	accountID, pageID, compID, watchdogID string,
	params *CreateWatchdogParams,
) (*Watchdog, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		watchdogPath(accountID, pageID, compID)+"/"+watchdogID, params,
	)
	if err != nil {
		return nil, err
	}
	var w Watchdog
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("unmarshal watchdog: %w", err)
	}
	return &w, nil
}

// DeleteWatchdog deletes a watchdog by ID.
func (c *Client) DeleteWatchdog(
	ctx context.Context, accountID, pageID, compID, watchdogID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		watchdogPath(accountID, pageID, compID)+"/"+watchdogID, nil,
	)
	return err
}
