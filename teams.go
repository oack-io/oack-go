package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Team represents an Oack team.
type Team struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TeamMember represents a user's membership in a team.
type TeamMember struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
}

// TeamInvite represents a pending team invitation.
type TeamInvite struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Token     string `json:"token"`
	Role      string `json:"role"`
	CreatedBy string `json:"created_by"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

// AcceptInviteResult is returned after successfully accepting an invite.
type AcceptInviteResult struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Role     string `json:"role"`
}

// TeamAPIKey represents an API key belonging to a team.
type TeamAPIKey struct {
	ID        string  `json:"id"`
	TeamID    string  `json:"team_id"`
	Name      string  `json:"name"`
	KeyPrefix string  `json:"key_prefix"`
	CreatedBy string  `json:"created_by"`
	ExpiresAt *string `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
}

// CreateTeamAPIKeyParams holds parameters for creating a team API key.
type CreateTeamAPIKeyParams struct {
	Name      string  `json:"name"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

// CreateTeamAPIKeyResult is returned after creating a team API key.
type CreateTeamAPIKeyResult struct {
	Key    string     `json:"key"`
	APIKey TeamAPIKey `json:"api_key"`
}

// CreateTeam creates a new team under the given account.
func (c *Client) CreateTeam(ctx context.Context, accountID, name string) (*Team, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/accounts/"+accountID+"/teams", map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	var t Team
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &t, nil
}

// ListTeams returns all teams accessible to the authenticated user.
func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return nil, err
	}
	var teams []Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, fmt.Errorf("unmarshal teams: %w", err)
	}
	return teams, nil
}

// ListAccountTeams returns all teams belonging to the given account.
func (c *Client) ListAccountTeams(ctx context.Context, accountID string) ([]Team, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/accounts/"+accountID+"/teams", nil)
	if err != nil {
		return nil, err
	}
	var teams []Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return nil, fmt.Errorf("unmarshal teams: %w", err)
	}
	return teams, nil
}

// GetTeam returns a single team by ID.
func (c *Client) GetTeam(ctx context.Context, teamID string) (*Team, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams/"+teamID, nil)
	if err != nil {
		return nil, err
	}
	var t Team
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &t, nil
}

// UpdateTeam updates the name of an existing team.
func (c *Client) UpdateTeam(ctx context.Context, teamID, name string) (*Team, error) {
	body, err := c.do(ctx, http.MethodPut, "/api/v1/teams/"+teamID, map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	var t Team
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("unmarshal team: %w", err)
	}
	return &t, nil
}

// DeleteTeam deletes a team by ID.
func (c *Client) DeleteTeam(ctx context.Context, teamID string) error {
	_, err := c.do(ctx, http.MethodDelete, "/api/v1/teams/"+teamID, nil)
	return err
}

// ListMembers returns all members of a team.
func (c *Client) ListMembers(ctx context.Context, teamID string) ([]TeamMember, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams/"+teamID+"/members", nil)
	if err != nil {
		return nil, err
	}
	var members []TeamMember
	if err := json.Unmarshal(body, &members); err != nil {
		return nil, fmt.Errorf("unmarshal members: %w", err)
	}
	return members, nil
}

// AddMember adds a user to a team with the specified role.
func (c *Client) AddMember(
	ctx context.Context, teamID, userID, role string,
) (*TeamMember, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/teams/"+teamID+"/members", map[string]string{
		"user_id": userID,
		"role":    role,
	})
	if err != nil {
		return nil, err
	}
	var m TeamMember
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal member: %w", err)
	}
	return &m, nil
}

// RemoveMember removes a user from a team.
func (c *Client) RemoveMember(ctx context.Context, teamID, userID string) error {
	_, err := c.do(ctx, http.MethodDelete, "/api/v1/teams/"+teamID+"/members/"+userID, nil)
	return err
}

// SetMemberRole updates a team member's role.
func (c *Client) SetMemberRole(ctx context.Context, teamID, userID, role string) error {
	_, err := c.do(
		ctx, http.MethodPut,
		"/api/v1/teams/"+teamID+"/members/"+userID+"/role",
		map[string]string{"role": role},
	)
	return err
}

// ListInvites returns all pending invites for a team.
func (c *Client) ListInvites(ctx context.Context, teamID string) ([]TeamInvite, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams/"+teamID+"/invites", nil)
	if err != nil {
		return nil, err
	}
	var invites []TeamInvite
	if err := json.Unmarshal(body, &invites); err != nil {
		return nil, fmt.Errorf("unmarshal invites: %w", err)
	}
	return invites, nil
}

// CreateInvite creates a new invite link for a team.
func (c *Client) CreateInvite(ctx context.Context, teamID string) (*TeamInvite, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/teams/"+teamID+"/invites", nil)
	if err != nil {
		return nil, err
	}
	var inv TeamInvite
	if err := json.Unmarshal(body, &inv); err != nil {
		return nil, fmt.Errorf("unmarshal invite: %w", err)
	}
	return &inv, nil
}

// RevokeInvite revokes a pending team invite.
func (c *Client) RevokeInvite(ctx context.Context, teamID, inviteID string) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		"/api/v1/teams/"+teamID+"/invites/"+inviteID, nil,
	)
	return err
}

// AcceptInvite accepts a team invitation using the invite token.
func (c *Client) AcceptInvite(ctx context.Context, token string) (*AcceptInviteResult, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/invites/"+token+"/accept", nil)
	if err != nil {
		return nil, err
	}
	var result AcceptInviteResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal accept invite result: %w", err)
	}
	return &result, nil
}

// CreateTeamAPIKey creates a new API key for a team.
func (c *Client) CreateTeamAPIKey(
	ctx context.Context, teamID string, params *CreateTeamAPIKeyParams,
) (*CreateTeamAPIKeyResult, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/teams/"+teamID+"/api-keys", params)
	if err != nil {
		return nil, err
	}
	var result CreateTeamAPIKeyResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal api key result: %w", err)
	}
	return &result, nil
}

// ListTeamAPIKeys returns all API keys for a team.
func (c *Client) ListTeamAPIKeys(ctx context.Context, teamID string) ([]TeamAPIKey, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams/"+teamID+"/api-keys", nil)
	if err != nil {
		return nil, err
	}
	var keys []TeamAPIKey
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, fmt.Errorf("unmarshal api keys: %w", err)
	}
	return keys, nil
}

// DeleteTeamAPIKey deletes an API key from a team.
func (c *Client) DeleteTeamAPIKey(ctx context.Context, teamID, keyID string) error {
	_, err := c.do(ctx, http.MethodDelete, "/api/v1/teams/"+teamID+"/api-keys/"+keyID, nil)
	return err
}
