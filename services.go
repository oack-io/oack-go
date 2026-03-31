package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Service represents a service in the incident management system.
type Service struct {
	ID                 string   `json:"id"`
	AccountID          string   `json:"account_id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	IntegrationKey     string   `json:"integration_key"`
	EscalationPolicyID string   `json:"escalation_policy_id,omitempty"`
	Status             string   `json:"status"`
	Tags               []string `json:"tags"`
	MonitorIDs         []string `json:"monitor_ids,omitempty"`
	ComponentIDs       []string `json:"component_ids,omitempty"`
	DependencyIDs      []string `json:"dependency_ids,omitempty"`
	DependentIDs       []string `json:"dependent_ids,omitempty"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

// CreateServiceParams holds parameters for creating a service.
type CreateServiceParams struct {
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	EscalationPolicyID string   `json:"escalation_policy_id,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	MonitorIDs         []string `json:"monitor_ids,omitempty"`
}

// UpdateServiceParams holds parameters for updating a service.
type UpdateServiceParams struct {
	Name               *string   `json:"name,omitempty"`
	Description        *string   `json:"description,omitempty"`
	EscalationPolicyID *string   `json:"escalation_policy_id,omitempty"`
	Tags               *[]string `json:"tags,omitempty"`
}

// ServiceAnalytics holds incident analytics for a service.
type ServiceAnalytics struct {
	MTTR            *float64            `json:"mttr_seconds"`
	MTTF            *float64            `json:"mttf_seconds"`
	IncidentCount   int                 `json:"incident_count"`
	BySeverity      map[string]int      `json:"by_severity"`
	OpenActionItems int                 `json:"open_action_items"`
	MTTRBySeverity  map[string]*float64 `json:"mttr_by_severity"`
	UptimePct       *float64            `json:"uptime_pct,omitempty"`
}

// serviceBasePath returns the base URL path for service operations.
func serviceBasePath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/services"
}

// servicePath returns the URL path for a specific service.
func servicePath(accountID, serviceID string) string {
	return serviceBasePath(accountID) + "/" + serviceID
}

// CreateService creates a new service for an account.
func (c *Client) CreateService(
	ctx context.Context, accountID string, params *CreateServiceParams,
) (*Service, error) {
	body, err := c.do(ctx, http.MethodPost, serviceBasePath(accountID), params)
	if err != nil {
		return nil, err
	}
	var s Service
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal service: %w", err)
	}
	return &s, nil
}

// ListServices returns all services for an account.
func (c *Client) ListServices(ctx context.Context, accountID string) ([]Service, error) {
	body, err := c.do(ctx, http.MethodGet, serviceBasePath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var services []Service
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("unmarshal services: %w", err)
	}
	return services, nil
}

// GetService returns a single service by ID.
func (c *Client) GetService(
	ctx context.Context, accountID, serviceID string,
) (*Service, error) {
	body, err := c.do(ctx, http.MethodGet, servicePath(accountID, serviceID), nil)
	if err != nil {
		return nil, err
	}
	var s Service
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal service: %w", err)
	}
	return &s, nil
}

// UpdateService updates an existing service.
func (c *Client) UpdateService(
	ctx context.Context, accountID, serviceID string, params *UpdateServiceParams,
) (*Service, error) {
	body, err := c.do(ctx, http.MethodPut, servicePath(accountID, serviceID), params)
	if err != nil {
		return nil, err
	}
	var s Service
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal service: %w", err)
	}
	return &s, nil
}

// DeleteService deletes a service by ID.
func (c *Client) DeleteService(ctx context.Context, accountID, serviceID string) error {
	_, err := c.do(ctx, http.MethodDelete, servicePath(accountID, serviceID), nil)
	return err
}

// LinkServiceMonitors links monitors to a service.
func (c *Client) LinkServiceMonitors(
	ctx context.Context, accountID, serviceID string, monitorIDs []string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		servicePath(accountID, serviceID)+"/monitors",
		map[string][]string{"monitor_ids": monitorIDs},
	)
	return err
}

// UnlinkServiceMonitor unlinks a monitor from a service.
func (c *Client) UnlinkServiceMonitor(
	ctx context.Context, accountID, serviceID, monitorID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		servicePath(accountID, serviceID)+"/monitors/"+monitorID, nil,
	)
	return err
}

// LinkServiceIncidents links incidents to a service.
func (c *Client) LinkServiceIncidents(
	ctx context.Context, accountID, serviceID string, incidentIDs []string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		servicePath(accountID, serviceID)+"/incidents",
		map[string][]string{"incident_ids": incidentIDs},
	)
	return err
}

// UnlinkServiceIncident unlinks an incident from a service.
func (c *Client) UnlinkServiceIncident(
	ctx context.Context, accountID, serviceID, incidentID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		servicePath(accountID, serviceID)+"/incidents/"+incidentID, nil,
	)
	return err
}

// GetServiceAnalytics returns incident analytics for a service.
// Optional query params from and to (RFC3339) control the time range; defaults to last 30 days.
func (c *Client) GetServiceAnalytics(
	ctx context.Context, accountID, serviceID string,
) (*ServiceAnalytics, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		servicePath(accountID, serviceID)+"/analytics", nil,
	)
	if err != nil {
		return nil, err
	}
	var a ServiceAnalytics
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal service analytics: %w", err)
	}
	return &a, nil
}
