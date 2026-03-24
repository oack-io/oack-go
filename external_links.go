package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ExternalLink represents a configurable link template attached to monitors.
type ExternalLink struct {
	ID                string `json:"id"`
	TeamID            string `json:"team_id"`
	Name              string `json:"name"`
	URLTemplate       string `json:"url_template"`
	IconURL           string `json:"icon_url"`
	TimeWindowMinutes int    `json:"time_window_minutes"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// CreateExternalLinkParams holds parameters for creating or updating an external link.
type CreateExternalLinkParams struct {
	Name              string `json:"name"`
	URLTemplate       string `json:"url_template"`
	IconURL           string `json:"icon_url,omitempty"`
	TimeWindowMinutes int    `json:"time_window_minutes"`
}

// CreateExternalLink creates a new external link for a team.
func (c *Client) CreateExternalLink(
	ctx context.Context, teamID string, params *CreateExternalLinkParams,
) (*ExternalLink, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/external-links", params,
	)
	if err != nil {
		return nil, err
	}
	var link ExternalLink
	if err := json.Unmarshal(body, &link); err != nil {
		return nil, fmt.Errorf("unmarshal external link: %w", err)
	}
	return &link, nil
}

// ListExternalLinks returns all external links for a team.
func (c *Client) ListExternalLinks(
	ctx context.Context, teamID string,
) ([]ExternalLink, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/external-links", nil,
	)
	if err != nil {
		return nil, err
	}
	var links []ExternalLink
	if err := json.Unmarshal(body, &links); err != nil {
		return nil, fmt.Errorf("unmarshal external links: %w", err)
	}
	return links, nil
}

// GetExternalLink returns a single external link by ID.
func (c *Client) GetExternalLink(
	ctx context.Context, teamID, linkID string,
) (*ExternalLink, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/external-links/"+linkID, nil,
	)
	if err != nil {
		return nil, err
	}
	var link ExternalLink
	if err := json.Unmarshal(body, &link); err != nil {
		return nil, fmt.Errorf("unmarshal external link: %w", err)
	}
	return &link, nil
}

// UpdateExternalLink updates an existing external link.
func (c *Client) UpdateExternalLink(
	ctx context.Context, teamID, linkID string, params *CreateExternalLinkParams,
) (*ExternalLink, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		"/api/v1/teams/"+teamID+"/external-links/"+linkID, params,
	)
	if err != nil {
		return nil, err
	}
	var link ExternalLink
	if err := json.Unmarshal(body, &link); err != nil {
		return nil, fmt.Errorf("unmarshal external link: %w", err)
	}
	return &link, nil
}

// DeleteExternalLink deletes an external link by ID.
func (c *Client) DeleteExternalLink(ctx context.Context, teamID, linkID string) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/external-links/"+linkID, nil,
	)
	return err
}

// AssignExternalLink assigns an external link to a monitor.
func (c *Client) AssignExternalLink(
	ctx context.Context, teamID, monitorID, linkID string,
) error {
	_, err := c.do(
		ctx, http.MethodPost,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/external-links/"+linkID, nil,
	)
	return err
}

// UnassignExternalLink removes an external link assignment from a monitor.
func (c *Client) UnassignExternalLink(
	ctx context.Context, teamID, monitorID, linkID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/external-links/"+linkID, nil,
	)
	return err
}

// ListMonitorExternalLinks returns all external links assigned to a monitor.
func (c *Client) ListMonitorExternalLinks(
	ctx context.Context, teamID, monitorID string,
) ([]ExternalLink, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/teams/"+teamID+"/monitors/"+monitorID+"/external-links", nil,
	)
	if err != nil {
		return nil, err
	}
	var links []ExternalLink
	if err := json.Unmarshal(body, &links); err != nil {
		return nil, fmt.Errorf("unmarshal external links: %w", err)
	}
	return links, nil
}
