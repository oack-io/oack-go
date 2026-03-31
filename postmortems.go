package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Postmortem is a structured retrospective for an incident.
type Postmortem struct {
	ID          string             `json:"id"`
	IncidentID  string             `json:"incident_id"`
	AccountID   string             `json:"account_id"`
	Status      string             `json:"status"`
	Summary     string             `json:"summary"`
	TimelineMD  string             `json:"timeline_md"`
	RootCauseMD string             `json:"root_cause_md"`
	ImpactMD    string             `json:"impact_md"`
	LessonsMD   string             `json:"lessons_md"`
	BodyMD      string             `json:"body_md"`
	ShareToken  string             `json:"share_token,omitempty"`
	CreatedBy   string             `json:"created_by"`
	PublishedAt string             `json:"published_at,omitempty"`
	ActionItems []PostmortemAction `json:"action_items"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

// PostmortemAction is a tracked action item from a postmortem.
type PostmortemAction struct {
	ID           string `json:"id"`
	PostmortemID string `json:"postmortem_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	OwnerID      string `json:"owner_id,omitempty"`
	Status       string `json:"status"`
	Priority     string `json:"priority"`
	DueDate      string `json:"due_date,omitempty"`
	CompletedAt  string `json:"completed_at,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreatePostmortemParams holds parameters for creating a postmortem.
type CreatePostmortemParams struct {
	BodyMD string `json:"body_md,omitempty"`
}

// UpdatePostmortemParams holds parameters for updating a postmortem.
type UpdatePostmortemParams struct {
	Summary     *string `json:"summary,omitempty"`
	TimelineMD  *string `json:"timeline_md,omitempty"`
	RootCauseMD *string `json:"root_cause_md,omitempty"`
	ImpactMD    *string `json:"impact_md,omitempty"`
	LessonsMD   *string `json:"lessons_md,omitempty"`
	BodyMD      *string `json:"body_md,omitempty"`
}

// CreateActionItemParams holds parameters for creating a postmortem action item.
type CreateActionItemParams struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	OwnerID     string `json:"owner_id,omitempty"`
	Priority    string `json:"priority,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
}

// UpdateActionItemParams holds parameters for updating a postmortem action item.
type UpdateActionItemParams struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	OwnerID     *string `json:"owner_id,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

// PostmortemTemplate is a reusable markdown template for postmortems.
type PostmortemTemplate struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsDefault bool   `json:"is_default"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreatePostmortemTemplateParams holds parameters for creating a postmortem template.
type CreatePostmortemTemplateParams struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsDefault bool   `json:"is_default"`
}

// UpdatePostmortemTemplateParams holds parameters for updating a postmortem template.
type UpdatePostmortemTemplateParams struct {
	Name      *string `json:"name,omitempty"`
	Content   *string `json:"content,omitempty"`
	IsDefault *bool   `json:"is_default,omitempty"`
}

func postmortemPath(accountID, incidentID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/incidents/%s/postmortem", accountID, incidentID)
}

func postmortemTemplatePath(accountID string) string {
	return fmt.Sprintf("/api/v1/accounts/%s/postmortem-templates", accountID)
}

// ── Postmortem CRUD

// CreatePostmortem creates a postmortem for an incident.
func (c *Client) CreatePostmortem(
	ctx context.Context, accountID, incidentID string, params *CreatePostmortemParams,
) (*Postmortem, error) {
	body, err := c.do(ctx, http.MethodPost, postmortemPath(accountID, incidentID), params)
	if err != nil {
		return nil, err
	}
	var pm Postmortem
	if err := json.Unmarshal(body, &pm); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem: %w", err)
	}
	return &pm, nil
}

// GetPostmortem returns the postmortem for an incident.
func (c *Client) GetPostmortem(
	ctx context.Context, accountID, incidentID string,
) (*Postmortem, error) {
	body, err := c.do(ctx, http.MethodGet, postmortemPath(accountID, incidentID), nil)
	if err != nil {
		return nil, err
	}
	var pm Postmortem
	if err := json.Unmarshal(body, &pm); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem: %w", err)
	}
	return &pm, nil
}

// UpdatePostmortem updates a postmortem.
func (c *Client) UpdatePostmortem(
	ctx context.Context, accountID, incidentID string, params *UpdatePostmortemParams,
) (*Postmortem, error) {
	body, err := c.do(ctx, http.MethodPut, postmortemPath(accountID, incidentID), params)
	if err != nil {
		return nil, err
	}
	var pm Postmortem
	if err := json.Unmarshal(body, &pm); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem: %w", err)
	}
	return &pm, nil
}

// DeletePostmortem deletes a postmortem.
func (c *Client) DeletePostmortem(
	ctx context.Context, accountID, incidentID string,
) error {
	_, err := c.do(ctx, http.MethodDelete, postmortemPath(accountID, incidentID), nil)
	return err
}

// PublishPostmortem publishes a postmortem.
func (c *Client) PublishPostmortem(
	ctx context.Context, accountID, incidentID string,
) (*Postmortem, error) {
	body, err := c.do(ctx, http.MethodPost, postmortemPath(accountID, incidentID)+"/publish", nil)
	if err != nil {
		return nil, err
	}
	var pm Postmortem
	if err := json.Unmarshal(body, &pm); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem: %w", err)
	}
	return &pm, nil
}

// GeneratePostmortemShareToken creates a shareable URL token.
func (c *Client) GeneratePostmortemShareToken(
	ctx context.Context, accountID, incidentID string,
) (string, error) {
	body, err := c.do(ctx, http.MethodPost, postmortemPath(accountID, incidentID)+"/share", nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		ShareToken string `json:"share_token"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("unmarshal share token: %w", err)
	}
	return resp.ShareToken, nil
}

// ── Action Items

// CreateActionItem adds an action item to a postmortem.
func (c *Client) CreateActionItem(
	ctx context.Context, accountID, incidentID string, params *CreateActionItemParams,
) (*PostmortemAction, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		postmortemPath(accountID, incidentID)+"/action-items", params,
	)
	if err != nil {
		return nil, err
	}
	var item PostmortemAction
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("unmarshal action item: %w", err)
	}
	return &item, nil
}

// UpdateActionItem updates an action item.
func (c *Client) UpdateActionItem(
	ctx context.Context, accountID, incidentID, itemID string, params *UpdateActionItemParams,
) (*PostmortemAction, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		postmortemPath(accountID, incidentID)+"/action-items/"+itemID, params,
	)
	if err != nil {
		return nil, err
	}
	var item PostmortemAction
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("unmarshal action item: %w", err)
	}
	return &item, nil
}

// DeleteActionItem deletes an action item.
func (c *Client) DeleteActionItem(
	ctx context.Context, accountID, incidentID, itemID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		postmortemPath(accountID, incidentID)+"/action-items/"+itemID, nil,
	)
	return err
}

// ── Postmortem Templates

// CreatePostmortemTemplate creates a new postmortem template.
func (c *Client) CreatePostmortemTemplate(
	ctx context.Context, accountID string, params *CreatePostmortemTemplateParams,
) (*PostmortemTemplate, error) {
	body, err := c.do(ctx, http.MethodPost, postmortemTemplatePath(accountID), params)
	if err != nil {
		return nil, err
	}
	var t PostmortemTemplate
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem template: %w", err)
	}
	return &t, nil
}

// ListPostmortemTemplates returns all postmortem templates for an account.
func (c *Client) ListPostmortemTemplates(
	ctx context.Context, accountID string,
) ([]PostmortemTemplate, error) {
	body, err := c.do(ctx, http.MethodGet, postmortemTemplatePath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var templates []PostmortemTemplate
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem templates: %w", err)
	}
	return templates, nil
}

// GetPostmortemTemplate returns a postmortem template by ID.
func (c *Client) GetPostmortemTemplate(
	ctx context.Context, accountID, templateID string,
) (*PostmortemTemplate, error) {
	body, err := c.do(ctx, http.MethodGet, postmortemTemplatePath(accountID)+"/"+templateID, nil)
	if err != nil {
		return nil, err
	}
	var t PostmortemTemplate
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem template: %w", err)
	}
	return &t, nil
}

// UpdatePostmortemTemplate updates a postmortem template.
func (c *Client) UpdatePostmortemTemplate(
	ctx context.Context, accountID, templateID string, params *UpdatePostmortemTemplateParams,
) (*PostmortemTemplate, error) {
	body, err := c.do(ctx, http.MethodPut, postmortemTemplatePath(accountID)+"/"+templateID, params)
	if err != nil {
		return nil, err
	}
	var t PostmortemTemplate
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal postmortem template: %w", err)
	}
	return &t, nil
}

// DeletePostmortemTemplate deletes a postmortem template.
func (c *Client) DeletePostmortemTemplate(
	ctx context.Context, accountID, templateID string,
) error {
	_, err := c.do(ctx, http.MethodDelete, postmortemTemplatePath(accountID)+"/"+templateID, nil)
	return err
}
