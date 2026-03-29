package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EnvVar represents a team-scoped environment variable or secret.
type EnvVar struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsSecret  bool   `json:"is_secret"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateEnvVarParams holds parameters for creating or updating an env var.
type CreateEnvVarParams struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}

// UpdateEnvVarParams holds parameters for updating an env var value.
type UpdateEnvVarParams struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}

func envBasePath(teamID string) string {
	return "/api/v1/teams/" + teamID + "/env"
}

// ListEnvVars returns all environment variables for a team.
// Secret values are masked in the response.
func (c *Client) ListEnvVars(ctx context.Context, teamID string) ([]EnvVar, error) {
	body, err := c.do(ctx, http.MethodGet, envBasePath(teamID), nil)
	if err != nil {
		return nil, err
	}
	var vars []EnvVar
	if err := json.Unmarshal(body, &vars); err != nil {
		return nil, fmt.Errorf("unmarshal env vars: %w", err)
	}
	return vars, nil
}

// CreateEnvVar creates a new environment variable or secret for a team.
func (c *Client) CreateEnvVar(
	ctx context.Context, teamID string, params *CreateEnvVarParams,
) (*EnvVar, error) {
	body, err := c.do(ctx, http.MethodPost, envBasePath(teamID), params)
	if err != nil {
		return nil, err
	}
	var v EnvVar
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("unmarshal env var: %w", err)
	}
	return &v, nil
}

// UpdateEnvVar updates an existing environment variable identified by key.
func (c *Client) UpdateEnvVar(
	ctx context.Context, teamID, key string, params *UpdateEnvVarParams,
) (*EnvVar, error) {
	body, err := c.do(ctx, http.MethodPut, envBasePath(teamID)+"/"+key, params)
	if err != nil {
		return nil, err
	}
	var v EnvVar
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("unmarshal env var: %w", err)
	}
	return &v, nil
}

// DeleteEnvVar deletes an environment variable by key.
func (c *Client) DeleteEnvVar(ctx context.Context, teamID, key string) error {
	_, err := c.do(ctx, http.MethodDelete, envBasePath(teamID)+"/"+key, nil)
	return err
}
