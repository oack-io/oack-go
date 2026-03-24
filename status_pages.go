package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// StatusPage represents a public status page.
type StatusPage struct {
	ID                   string  `json:"id"`
	AccountID            string  `json:"account_id"`
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug"`
	Description          string  `json:"description"`
	CustomDomain         *string `json:"custom_domain"`
	HasPassword          bool    `json:"has_password"`
	AllowIframe          bool    `json:"allow_iframe"`
	ShowHistoricalUptime bool    `json:"show_historical_uptime"`
	BrandingLogoURL      *string `json:"branding_logo_url"`
	BrandingFaviconURL   *string `json:"branding_favicon_url"`
	BrandingPrimaryColor *string `json:"branding_primary_color"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// CreateStatusPageParams holds parameters for creating or updating a status page.
type CreateStatusPageParams struct {
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug"`
	Description          string  `json:"description,omitempty"`
	CustomDomain         *string `json:"custom_domain,omitempty"`
	Password             *string `json:"password,omitempty"`
	AllowIframe          *bool   `json:"allow_iframe,omitempty"`
	ShowHistoricalUptime *bool   `json:"show_historical_uptime,omitempty"`
	BrandingLogoURL      *string `json:"branding_logo_url,omitempty"`
	BrandingFaviconURL   *string `json:"branding_favicon_url,omitempty"`
	BrandingPrimaryColor *string `json:"branding_primary_color,omitempty"`
}

// ComponentGroup represents a group of components on a status page.
type ComponentGroup struct {
	ID           string `json:"id"`
	StatusPageID string `json:"status_page_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Position     int    `json:"position"`
	Collapsed    bool   `json:"collapsed"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateComponentGroupParams holds parameters for creating or updating a component group.
type CreateComponentGroupParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Position    int    `json:"position"`
	Collapsed   *bool  `json:"collapsed,omitempty"`
}

// Component represents a single component on a status page.
type Component struct {
	ID            string `json:"id"`
	StatusPageID  string `json:"status_page_id"`
	GroupID       string `json:"group_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	DisplayUptime bool   `json:"display_uptime"`
	Position      int    `json:"position"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateComponentParams holds parameters for creating or updating a component.
type CreateComponentParams struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	GroupID       string `json:"group_id,omitempty"`
	DisplayUptime *bool  `json:"display_uptime,omitempty"`
	Position      int    `json:"position"`
}

// Incident represents a status page incident.
type Incident struct {
	ID           string  `json:"id"`
	StatusPageID string  `json:"status_page_id"`
	Title        string  `json:"title"`
	Message      string  `json:"message"`
	Severity     string  `json:"severity"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ResolvedAt   *string `json:"resolved_at"`
}

// CreateIncidentParams holds parameters for creating an incident.
type CreateIncidentParams struct {
	Title        string   `json:"title"`
	Message      string   `json:"message"`
	Severity     string   `json:"severity"`
	Status       string   `json:"status,omitempty"`
	ComponentIDs []string `json:"component_ids,omitempty"`
}

// UpdateIncidentParams holds parameters for updating an incident.
type UpdateIncidentParams struct {
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`
	Status   string `json:"status,omitempty"`
}

// IncidentUpdate represents an update posted to an incident.
type IncidentUpdate struct {
	ID         string `json:"id"`
	IncidentID string `json:"incident_id"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// Maintenance represents a scheduled maintenance window.
type Maintenance struct {
	ID           string  `json:"id"`
	StatusPageID string  `json:"status_page_id"`
	Title        string  `json:"title"`
	Message      string  `json:"message"`
	Status       string  `json:"status"`
	ScheduledAt  string  `json:"scheduled_at"`
	CompletedAt  *string `json:"completed_at"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// CreateMaintenanceParams holds parameters for creating a maintenance window.
type CreateMaintenanceParams struct {
	Title        string   `json:"title"`
	Message      string   `json:"message"`
	ScheduledAt  string   `json:"scheduled_at"`
	ComponentIDs []string `json:"component_ids,omitempty"`
}

// UpdateMaintenanceParams holds parameters for updating a maintenance window.
type UpdateMaintenanceParams struct {
	Title       string `json:"title,omitempty"`
	Message     string `json:"message,omitempty"`
	Status      string `json:"status,omitempty"`
	ScheduledAt string `json:"scheduled_at,omitempty"`
}

// MaintenanceUpdate represents an update posted to a maintenance window.
type MaintenanceUpdate struct {
	ID            string `json:"id"`
	MaintenanceID string `json:"maintenance_id"`
	Message       string `json:"message"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// Subscriber represents a status page subscriber.
type Subscriber struct {
	ID            string  `json:"id"`
	StatusPageID  string  `json:"status_page_id"`
	Email         *string `json:"email,omitempty"`
	EmailVerified bool    `json:"email_verified"`
	Scope         string  `json:"scope"`
	CreatedAt     string  `json:"created_at"`
	ConfirmedAt   *string `json:"confirmed_at,omitempty"`
}

// IncidentTemplate represents a reusable incident template.
type IncidentTemplate struct {
	ID           string `json:"id"`
	StatusPageID string `json:"status_page_id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Severity     string `json:"severity"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateIncidentTemplateParams holds parameters for creating an incident template.
type CreateIncidentTemplateParams struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// UpdateIncidentTemplateParams holds parameters for updating an incident template.
type UpdateIncidentTemplateParams struct {
	Name     *string `json:"name,omitempty"`
	Message  *string `json:"message,omitempty"`
	Severity *string `json:"severity,omitempty"`
}

func statusPagePath(accountID, pageID string) string {
	return "/api/v1/accounts/" + accountID + "/status-pages/" + pageID
}

// CreateStatusPage creates a new status page.
func (c *Client) CreateStatusPage(
	ctx context.Context, accountID string, params *CreateStatusPageParams,
) (*StatusPage, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		"/api/v1/accounts/"+accountID+"/status-pages", params,
	)
	if err != nil {
		return nil, err
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

// GetStatusPage returns a single status page by ID.
func (c *Client) GetStatusPage(
	ctx context.Context, accountID, pageID string,
) (*StatusPage, error) {
	body, err := c.do(ctx, http.MethodGet, statusPagePath(accountID, pageID), nil)
	if err != nil {
		return nil, err
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

// UpdateStatusPage updates an existing status page.
func (c *Client) UpdateStatusPage(
	ctx context.Context, accountID, pageID string, params *CreateStatusPageParams,
) (*StatusPage, error) {
	body, err := c.do(
		ctx, http.MethodPut, statusPagePath(accountID, pageID), params,
	)
	if err != nil {
		return nil, err
	}
	var sp StatusPage
	if err := json.Unmarshal(body, &sp); err != nil {
		return nil, fmt.Errorf("unmarshal status page: %w", err)
	}
	return &sp, nil
}

// DeleteStatusPage deletes a status page by ID.
func (c *Client) DeleteStatusPage(ctx context.Context, accountID, pageID string) error {
	_, err := c.do(ctx, http.MethodDelete, statusPagePath(accountID, pageID), nil)
	return err
}

// ListStatusPages returns all status pages for an account.
func (c *Client) ListStatusPages(ctx context.Context, accountID string) ([]StatusPage, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/accounts/"+accountID+"/status-pages", nil,
	)
	if err != nil {
		return nil, err
	}
	var pages []StatusPage
	if err := json.Unmarshal(body, &pages); err != nil {
		return nil, fmt.Errorf("unmarshal status pages: %w", err)
	}
	return pages, nil
}

// CreateComponentGroup creates a new component group on a status page.
func (c *Client) CreateComponentGroup(
	ctx context.Context, accountID, pageID string, params *CreateComponentGroupParams,
) (*ComponentGroup, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/component-groups", params,
	)
	if err != nil {
		return nil, err
	}
	var g ComponentGroup
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("unmarshal component group: %w", err)
	}
	return &g, nil
}

// GetComponentGroup returns a single component group by listing all groups and finding by ID.
func (c *Client) GetComponentGroup(
	ctx context.Context, accountID, pageID, groupID string,
) (*ComponentGroup, error) {
	groups, err := c.ListComponentGroups(ctx, accountID, pageID)
	if err != nil {
		return nil, err
	}
	for _, g := range groups {
		if g.ID == groupID {
			return &g, nil
		}
	}
	return nil, &APIError{StatusCode: http.StatusNotFound, Message: "component group not found"}
}

// UpdateComponentGroup updates an existing component group.
func (c *Client) UpdateComponentGroup(
	ctx context.Context,
	accountID, pageID, groupID string,
	params *CreateComponentGroupParams,
) (*ComponentGroup, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		statusPagePath(accountID, pageID)+"/component-groups/"+groupID, params,
	)
	if err != nil {
		return nil, err
	}
	var g ComponentGroup
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("unmarshal component group: %w", err)
	}
	return &g, nil
}

// DeleteComponentGroup deletes a component group by ID.
func (c *Client) DeleteComponentGroup(
	ctx context.Context, accountID, pageID, groupID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/component-groups/"+groupID, nil,
	)
	return err
}

// ListComponentGroups returns all component groups for a status page.
func (c *Client) ListComponentGroups(
	ctx context.Context, accountID, pageID string,
) ([]ComponentGroup, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/component-groups", nil,
	)
	if err != nil {
		return nil, err
	}
	var groups []ComponentGroup
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, fmt.Errorf("unmarshal component groups: %w", err)
	}
	return groups, nil
}

// CreateComponent creates a new component on a status page.
func (c *Client) CreateComponent(
	ctx context.Context, accountID, pageID string, params *CreateComponentParams,
) (*Component, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/components", params,
	)
	if err != nil {
		return nil, err
	}
	var comp Component
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, fmt.Errorf("unmarshal component: %w", err)
	}
	return &comp, nil
}

// GetComponent returns a single component by listing all components and finding by ID.
func (c *Client) GetComponent(
	ctx context.Context, accountID, pageID, compID string,
) (*Component, error) {
	components, err := c.ListComponents(ctx, accountID, pageID)
	if err != nil {
		return nil, err
	}
	for _, comp := range components {
		if comp.ID == compID {
			return &comp, nil
		}
	}
	return nil, &APIError{StatusCode: http.StatusNotFound, Message: "component not found"}
}

// UpdateComponent updates an existing component.
func (c *Client) UpdateComponent(
	ctx context.Context,
	accountID, pageID, compID string,
	params *CreateComponentParams,
) (*Component, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		statusPagePath(accountID, pageID)+"/components/"+compID, params,
	)
	if err != nil {
		return nil, err
	}
	var comp Component
	if err := json.Unmarshal(body, &comp); err != nil {
		return nil, fmt.Errorf("unmarshal component: %w", err)
	}
	return &comp, nil
}

// DeleteComponent deletes a component by ID.
func (c *Client) DeleteComponent(
	ctx context.Context, accountID, pageID, compID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/components/"+compID, nil,
	)
	return err
}

// ListComponents returns all components for a status page.
func (c *Client) ListComponents(
	ctx context.Context, accountID, pageID string,
) ([]Component, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/components", nil,
	)
	if err != nil {
		return nil, err
	}
	var components []Component
	if err := json.Unmarshal(body, &components); err != nil {
		return nil, fmt.Errorf("unmarshal components: %w", err)
	}
	return components, nil
}

// CreateIncident creates a new incident on a status page.
func (c *Client) CreateIncident(
	ctx context.Context, accountID, pageID string, params *CreateIncidentParams,
) (*Incident, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/incidents", params,
	)
	if err != nil {
		return nil, err
	}
	var inc Incident
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal incident: %w", err)
	}
	return &inc, nil
}

// GetIncident returns a single incident by ID.
func (c *Client) GetIncident(
	ctx context.Context, accountID, pageID, incidentID string,
) (*Incident, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/incidents/"+incidentID, nil,
	)
	if err != nil {
		return nil, err
	}
	var inc Incident
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal incident: %w", err)
	}
	return &inc, nil
}

// UpdateIncident updates an existing incident.
func (c *Client) UpdateIncident(
	ctx context.Context,
	accountID, pageID, incidentID string,
	params *UpdateIncidentParams,
) (*Incident, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		statusPagePath(accountID, pageID)+"/incidents/"+incidentID, params,
	)
	if err != nil {
		return nil, err
	}
	var inc Incident
	if err := json.Unmarshal(body, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal incident: %w", err)
	}
	return &inc, nil
}

// DeleteIncident deletes an incident by ID.
func (c *Client) DeleteIncident(
	ctx context.Context, accountID, pageID, incidentID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/incidents/"+incidentID, nil,
	)
	return err
}

// ListIncidents returns all incidents for a status page.
func (c *Client) ListIncidents(
	ctx context.Context, accountID, pageID string,
) ([]Incident, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/incidents", nil,
	)
	if err != nil {
		return nil, err
	}
	var incidents []Incident
	if err := json.Unmarshal(body, &incidents); err != nil {
		return nil, fmt.Errorf("unmarshal incidents: %w", err)
	}
	return incidents, nil
}

// PostIncidentUpdate posts a status update to an incident.
func (c *Client) PostIncidentUpdate(
	ctx context.Context, accountID, pageID, incidentID, message, status string,
) (*IncidentUpdate, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/incidents/"+incidentID+"/updates",
		map[string]string{"message": message, "status": status},
	)
	if err != nil {
		return nil, err
	}
	var u IncidentUpdate
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal incident update: %w", err)
	}
	return &u, nil
}

// CreateMaintenance creates a new maintenance window on a status page.
func (c *Client) CreateMaintenance(
	ctx context.Context, accountID, pageID string, params *CreateMaintenanceParams,
) (*Maintenance, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/maintenances", params,
	)
	if err != nil {
		return nil, err
	}
	var m Maintenance
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal maintenance: %w", err)
	}
	return &m, nil
}

// GetMaintenance returns a single maintenance window by ID.
func (c *Client) GetMaintenance(
	ctx context.Context, accountID, pageID, maintID string,
) (*Maintenance, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/maintenances/"+maintID, nil,
	)
	if err != nil {
		return nil, err
	}
	var m Maintenance
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal maintenance: %w", err)
	}
	return &m, nil
}

// UpdateMaintenance updates an existing maintenance window.
func (c *Client) UpdateMaintenance(
	ctx context.Context,
	accountID, pageID, maintID string,
	params *UpdateMaintenanceParams,
) (*Maintenance, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		statusPagePath(accountID, pageID)+"/maintenances/"+maintID, params,
	)
	if err != nil {
		return nil, err
	}
	var m Maintenance
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal maintenance: %w", err)
	}
	return &m, nil
}

// DeleteMaintenance deletes a maintenance window by ID.
func (c *Client) DeleteMaintenance(
	ctx context.Context, accountID, pageID, maintID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/maintenances/"+maintID, nil,
	)
	return err
}

// ListMaintenances returns all maintenance windows for a status page.
func (c *Client) ListMaintenances(
	ctx context.Context, accountID, pageID string,
) ([]Maintenance, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/maintenances", nil,
	)
	if err != nil {
		return nil, err
	}
	var maintenances []Maintenance
	if err := json.Unmarshal(body, &maintenances); err != nil {
		return nil, fmt.Errorf("unmarshal maintenances: %w", err)
	}
	return maintenances, nil
}

// PostMaintenanceUpdate posts a status update to a maintenance window.
func (c *Client) PostMaintenanceUpdate(
	ctx context.Context, accountID, pageID, maintID, message, status string,
) (*MaintenanceUpdate, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/maintenances/"+maintID+"/updates",
		map[string]string{"message": message, "status": status},
	)
	if err != nil {
		return nil, err
	}
	var u MaintenanceUpdate
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal maintenance update: %w", err)
	}
	return &u, nil
}

// ListSubscribers returns all subscribers for a status page.
func (c *Client) ListSubscribers(
	ctx context.Context, accountID, pageID string,
) ([]Subscriber, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/subscribers", nil,
	)
	if err != nil {
		return nil, err
	}
	var subs []Subscriber
	if err := json.Unmarshal(body, &subs); err != nil {
		return nil, fmt.Errorf("unmarshal subscribers: %w", err)
	}
	return subs, nil
}

// RemoveSubscriber removes a subscriber from a status page.
func (c *Client) RemoveSubscriber(
	ctx context.Context, accountID, pageID, subscriberID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/subscribers/"+subscriberID, nil,
	)
	return err
}

// CreateIncidentTemplate creates a new incident template on a status page.
func (c *Client) CreateIncidentTemplate(
	ctx context.Context, accountID, pageID string, params *CreateIncidentTemplateParams,
) (*IncidentTemplate, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		statusPagePath(accountID, pageID)+"/incident-templates", params,
	)
	if err != nil {
		return nil, err
	}
	var t IncidentTemplate
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal incident template: %w", err)
	}
	return &t, nil
}

// ListIncidentTemplates returns all incident templates for a status page.
func (c *Client) ListIncidentTemplates(
	ctx context.Context, accountID, pageID string,
) ([]IncidentTemplate, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		statusPagePath(accountID, pageID)+"/incident-templates", nil,
	)
	if err != nil {
		return nil, err
	}
	var templates []IncidentTemplate
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("unmarshal incident templates: %w", err)
	}
	return templates, nil
}

// UpdateIncidentTemplate updates an existing incident template.
func (c *Client) UpdateIncidentTemplate(
	ctx context.Context,
	accountID, pageID, templateID string,
	params *UpdateIncidentTemplateParams,
) (*IncidentTemplate, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		statusPagePath(accountID, pageID)+"/incident-templates/"+templateID, params,
	)
	if err != nil {
		return nil, err
	}
	var t IncidentTemplate
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal incident template: %w", err)
	}
	return &t, nil
}

// DeleteIncidentTemplate deletes an incident template by ID.
func (c *Client) DeleteIncidentTemplate(
	ctx context.Context, accountID, pageID, templateID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		statusPagePath(accountID, pageID)+"/incident-templates/"+templateID, nil,
	)
	return err
}
