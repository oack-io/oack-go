package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PDIntegration represents a PagerDuty integration.
type PDIntegration struct {
	ID           string   `json:"id"`
	AccountID    string   `json:"account_id"`
	APIKey       string   `json:"api_key"`
	Region       string   `json:"region"`
	ServiceIDs   []string `json:"service_ids"`
	SyncEnabled  bool     `json:"sync_enabled"`
	SyncError    string   `json:"sync_error"`
	LastSyncedAt *string  `json:"last_synced_at"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// CreatePDIntegrationParams holds parameters for creating a PagerDuty integration.
type CreatePDIntegrationParams struct {
	APIKey     string   `json:"api_key"`
	Region     string   `json:"region"`
	ServiceIDs []string `json:"service_ids,omitempty"`
}

// UpdatePDIntegrationParams holds parameters for updating a PagerDuty integration.
type UpdatePDIntegrationParams struct {
	APIKey      *string  `json:"api_key,omitempty"`
	Region      *string  `json:"region,omitempty"`
	ServiceIDs  []string `json:"service_ids,omitempty"`
	SyncEnabled *bool    `json:"sync_enabled,omitempty"`
}

// CFIntegration represents a Cloudflare zone integration.
type CFIntegration struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	ZoneID          string  `json:"zone_id"`
	ZoneName        string  `json:"zone_name"`
	APIToken        string  `json:"api_token"`
	Enabled         bool    `json:"enabled"`
	SessionError    string  `json:"session_error"`
	LastConnectedAt *string `json:"last_connected_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// CreateCFIntegrationParams holds parameters for creating a Cloudflare zone integration.
type CreateCFIntegrationParams struct {
	ZoneID   string `json:"zone_id"`
	ZoneName string `json:"zone_name"`
	APIToken string `json:"api_token"`
}

func pdPath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/integrations/pagerduty"
}

func cfPath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/integrations/cloudflare-zone"
}

// CreatePDIntegration creates a new PagerDuty integration.
func (c *Client) CreatePDIntegration(
	ctx context.Context, accountID string, params *CreatePDIntegrationParams,
) (*PDIntegration, error) {
	body, err := c.do(ctx, http.MethodPost, pdPath(accountID), params)
	if err != nil {
		return nil, err
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pagerduty integration: %w", err)
	}
	return &pd, nil
}

// GetPDIntegration returns the PagerDuty integration for an account.
func (c *Client) GetPDIntegration(
	ctx context.Context, accountID string,
) (*PDIntegration, error) {
	body, err := c.do(ctx, http.MethodGet, pdPath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pagerduty integration: %w", err)
	}
	return &pd, nil
}

// UpdatePDIntegration updates the PagerDuty integration for an account.
func (c *Client) UpdatePDIntegration(
	ctx context.Context, accountID string, params *UpdatePDIntegrationParams,
) (*PDIntegration, error) {
	body, err := c.do(ctx, http.MethodPut, pdPath(accountID), params)
	if err != nil {
		return nil, err
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pagerduty integration: %w", err)
	}
	return &pd, nil
}

// DeletePDIntegration deletes the PagerDuty integration for an account.
func (c *Client) DeletePDIntegration(ctx context.Context, accountID string) error {
	_, err := c.do(ctx, http.MethodDelete, pdPath(accountID), nil)
	return err
}

// SyncPDIntegration triggers a sync of the PagerDuty integration.
func (c *Client) SyncPDIntegration(
	ctx context.Context, accountID string,
) (*PDIntegration, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		pdPath(accountID)+"/sync", nil,
	)
	if err != nil {
		return nil, err
	}
	var pd PDIntegration
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, fmt.Errorf("unmarshal pagerduty integration: %w", err)
	}
	return &pd, nil
}

// CreateCFIntegration creates a new Cloudflare zone integration.
func (c *Client) CreateCFIntegration(
	ctx context.Context, accountID string, params *CreateCFIntegrationParams,
) (*CFIntegration, error) {
	body, err := c.do(ctx, http.MethodPost, cfPath(accountID), params)
	if err != nil {
		return nil, err
	}
	var cf CFIntegration
	if err := json.Unmarshal(body, &cf); err != nil {
		return nil, fmt.Errorf("unmarshal cloudflare integration: %w", err)
	}
	return &cf, nil
}

// ListCFIntegrations returns all Cloudflare zone integrations for an account.
func (c *Client) ListCFIntegrations(
	ctx context.Context, accountID string,
) ([]CFIntegration, error) {
	body, err := c.do(ctx, http.MethodGet, cfPath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var integrations []CFIntegration
	if err := json.Unmarshal(body, &integrations); err != nil {
		return nil, fmt.Errorf("unmarshal cloudflare integrations: %w", err)
	}
	return integrations, nil
}

// UpdateCFIntegration updates the API token for a Cloudflare zone integration.
func (c *Client) UpdateCFIntegration(
	ctx context.Context, accountID, cfID, apiToken string,
) (*CFIntegration, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		cfPath(accountID)+"/"+cfID,
		map[string]string{"api_token": apiToken},
	)
	if err != nil {
		return nil, err
	}
	var cf CFIntegration
	if err := json.Unmarshal(body, &cf); err != nil {
		return nil, fmt.Errorf("unmarshal cloudflare integration: %w", err)
	}
	return &cf, nil
}

// DeleteCFIntegration deletes a Cloudflare zone integration by ID.
func (c *Client) DeleteCFIntegration(ctx context.Context, accountID, cfID string) error {
	_, err := c.do(ctx, http.MethodDelete, cfPath(accountID)+"/"+cfID, nil)
	return err
}
