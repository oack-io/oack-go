package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Trigger links a monitor to a status page component for automatic incident management.
type Trigger struct {
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

// CreateTriggerParams holds parameters for creating or updating a trigger.
type CreateTriggerParams struct {
	MonitorID         string `json:"monitor_id"`
	Severity          string `json:"severity"`
	AutoCreate        *bool  `json:"auto_create,omitempty"`
	AutoResolve       *bool  `json:"auto_resolve,omitempty"`
	NotifySubscribers *bool  `json:"notify_subscribers,omitempty"`
	TemplateID        string `json:"template_id,omitempty"`
}

func triggerPath(accountID, pageID, compID string) string {
	return statusPagePath(accountID, pageID) + "/components/" + compID + "/triggers"
}

// CreateTrigger creates a new trigger on a status page component.
func (c *Client) CreateTrigger(
	ctx context.Context, accountID, pageID, compID string, params *CreateTriggerParams,
) (*Trigger, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		triggerPath(accountID, pageID, compID), params,
	)
	if err != nil {
		return nil, err
	}
	var t Trigger
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal trigger: %w", err)
	}
	return &t, nil
}

// ListTriggers returns all triggers for a status page component.
func (c *Client) ListTriggers(
	ctx context.Context, accountID, pageID, compID string,
) ([]Trigger, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		triggerPath(accountID, pageID, compID), nil,
	)
	if err != nil {
		return nil, err
	}
	var triggers []Trigger
	if err := json.Unmarshal(body, &triggers); err != nil {
		return nil, fmt.Errorf("unmarshal triggers: %w", err)
	}
	return triggers, nil
}

// UpdateTrigger updates an existing trigger.
func (c *Client) UpdateTrigger(
	ctx context.Context,
	accountID, pageID, compID, triggerID string,
	params *CreateTriggerParams,
) (*Trigger, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		triggerPath(accountID, pageID, compID)+"/"+triggerID, params,
	)
	if err != nil {
		return nil, err
	}
	var t Trigger
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal trigger: %w", err)
	}
	return &t, nil
}

// DeleteTrigger deletes a trigger by ID.
func (c *Client) DeleteTrigger(
	ctx context.Context, accountID, pageID, compID, triggerID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		triggerPath(accountID, pageID, compID)+"/"+triggerID, nil,
	)
	return err
}
