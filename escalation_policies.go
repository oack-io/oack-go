package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EscalationPolicy defines how alerts escalate through on-call levels.
type EscalationPolicy struct {
	ID        string            `json:"id"`
	AccountID string            `json:"account_id"`
	Name      string            `json:"name"`
	Levels    []EscalationLevel `json:"levels"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

// EscalationLevel defines a single level in an escalation policy.
type EscalationLevel struct {
	ScheduleID        string `json:"schedule_id"`
	AckTimeoutMinutes int    `json:"ack_timeout_minutes,omitempty"`
	DurationMinutes   int    `json:"duration_minutes,omitempty"`
}

// CreateEscalationPolicyParams holds parameters for creating an escalation policy.
type CreateEscalationPolicyParams struct {
	Name   string            `json:"name"`
	Levels []EscalationLevel `json:"levels,omitempty"`
}

// UpdateEscalationPolicyParams holds parameters for updating an escalation policy.
type UpdateEscalationPolicyParams struct {
	Name   string            `json:"name,omitempty"`
	Levels []EscalationLevel `json:"levels,omitempty"`
}

// escalationBasePath returns the base URL path for escalation policy operations.
func escalationBasePath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/oncall/escalation-policies"
}

// escalationPath returns the URL path for a specific escalation policy.
func escalationPath(accountID, policyID string) string {
	return escalationBasePath(accountID) + "/" + policyID
}

// CreateEscalationPolicy creates a new escalation policy.
func (c *Client) CreateEscalationPolicy(
	ctx context.Context, accountID string, params *CreateEscalationPolicyParams,
) (*EscalationPolicy, error) {
	body, err := c.do(ctx, http.MethodPost, escalationBasePath(accountID), params)
	if err != nil {
		return nil, err
	}
	var p EscalationPolicy
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("unmarshal escalation policy: %w", err)
	}
	return &p, nil
}

// ListEscalationPolicies returns all escalation policies for an account.
func (c *Client) ListEscalationPolicies(
	ctx context.Context, accountID string,
) ([]EscalationPolicy, error) {
	body, err := c.do(ctx, http.MethodGet, escalationBasePath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var policies []EscalationPolicy
	if err := json.Unmarshal(body, &policies); err != nil {
		return nil, fmt.Errorf("unmarshal escalation policies: %w", err)
	}
	return policies, nil
}

// UpdateEscalationPolicy updates an existing escalation policy.
func (c *Client) UpdateEscalationPolicy(
	ctx context.Context, accountID, policyID string, params *UpdateEscalationPolicyParams,
) (*EscalationPolicy, error) {
	body, err := c.do(ctx, http.MethodPut, escalationPath(accountID, policyID), params)
	if err != nil {
		return nil, err
	}
	var p EscalationPolicy
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("unmarshal escalation policy: %w", err)
	}
	return &p, nil
}

// DeleteEscalationPolicy deletes an escalation policy by ID.
func (c *Client) DeleteEscalationPolicy(
	ctx context.Context, accountID, policyID string,
) error {
	_, err := c.do(ctx, http.MethodDelete, escalationPath(accountID, policyID), nil)
	return err
}
