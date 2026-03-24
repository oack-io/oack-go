package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Account represents an Oack account.
type Account struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Plan      string  `json:"plan"`
	DeletedAt *string `json:"deleted_at"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// AccountMember represents a user's membership in an account.
type AccountMember struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
}

// AccountInvite represents a pending account invitation.
type AccountInvite struct {
	ID         string  `json:"id"`
	AccountID  string  `json:"account_id"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	InvitedBy  string  `json:"invited_by"`
	Token      string  `json:"token"`
	InviteURL  string  `json:"invite_url"`
	ExpiresAt  string  `json:"expires_at"`
	AcceptedAt *string `json:"accepted_at,omitempty"`
	RevokedAt  *string `json:"revoked_at,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

// Subscription represents an account's subscription.
type Subscription struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id"`
	Plan      string  `json:"plan"`
	Status    string  `json:"status"`
	ExpiresAt *string `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// accountPath returns the URL path for a specific account.
func accountPath(accountID string) string {
	return "/api/v1/accounts/" + accountID
}

// CreateAccount creates a new account.
func (c *Client) CreateAccount(ctx context.Context, name string) (*Account, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/accounts", map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	var a Account
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &a, nil
}

// ListAccounts returns all accounts accessible to the authenticated user.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/accounts", nil)
	if err != nil {
		return nil, err
	}
	var accounts []Account
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, fmt.Errorf("unmarshal accounts: %w", err)
	}
	return accounts, nil
}

// GetAccount returns a single account by ID.
func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	body, err := c.do(ctx, http.MethodGet, accountPath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var a Account
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &a, nil
}

// UpdateAccount updates the name of an existing account.
func (c *Client) UpdateAccount(
	ctx context.Context, accountID, name string,
) (*Account, error) {
	body, err := c.do(ctx, http.MethodPut, accountPath(accountID), map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	var a Account
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &a, nil
}

// DeleteAccount deletes an account by ID.
func (c *Client) DeleteAccount(ctx context.Context, accountID string) error {
	_, err := c.do(ctx, http.MethodDelete, accountPath(accountID), nil)
	return err
}

// RestoreAccount restores a previously deleted account.
func (c *Client) RestoreAccount(ctx context.Context, accountID string) (*Account, error) {
	body, err := c.do(ctx, http.MethodPost, accountPath(accountID)+"/restore", nil)
	if err != nil {
		return nil, err
	}
	var a Account
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &a, nil
}

// TransferAccount transfers account ownership to another user.
func (c *Client) TransferAccount(
	ctx context.Context, accountID, userID string,
) (*Account, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		accountPath(accountID)+"/transfer",
		map[string]string{"user_id": userID},
	)
	if err != nil {
		return nil, err
	}
	var a Account
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("unmarshal account: %w", err)
	}
	return &a, nil
}

// ListAccountMembers returns all members of an account.
func (c *Client) ListAccountMembers(
	ctx context.Context, accountID string,
) ([]AccountMember, error) {
	body, err := c.do(ctx, http.MethodGet, accountPath(accountID)+"/members", nil)
	if err != nil {
		return nil, err
	}
	var members []AccountMember
	if err := json.Unmarshal(body, &members); err != nil {
		return nil, fmt.Errorf("unmarshal account members: %w", err)
	}
	return members, nil
}

// SetAccountMemberRole updates a member's role within an account.
func (c *Client) SetAccountMemberRole(
	ctx context.Context, accountID, userID, role string,
) (*AccountMember, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		accountPath(accountID)+"/members/"+userID+"/role",
		map[string]string{"role": role},
	)
	if err != nil {
		return nil, err
	}
	var m AccountMember
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("unmarshal account member: %w", err)
	}
	return &m, nil
}

// RemoveAccountMember removes a user from an account.
func (c *Client) RemoveAccountMember(ctx context.Context, accountID, userID string) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		accountPath(accountID)+"/members/"+userID, nil,
	)
	return err
}

// GetAccountSubscription returns the subscription for an account.
func (c *Client) GetAccountSubscription(
	ctx context.Context, accountID string,
) (*Subscription, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		accountPath(accountID)+"/subscription", nil,
	)
	if err != nil {
		return nil, err
	}
	var s Subscription
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal subscription: %w", err)
	}
	return &s, nil
}

// UpdateAccountSubscription updates the plan and status of an account's subscription.
func (c *Client) UpdateAccountSubscription(
	ctx context.Context, accountID, plan, status string,
) (*Subscription, error) {
	body, err := c.do(
		ctx, http.MethodPut,
		accountPath(accountID)+"/subscription",
		map[string]string{"plan": plan, "status": status},
	)
	if err != nil {
		return nil, err
	}
	var s Subscription
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal subscription: %w", err)
	}
	return &s, nil
}

// CreateAccountInvite creates a new invitation for an account.
func (c *Client) CreateAccountInvite(
	ctx context.Context, accountID, email, role string,
) (*AccountInvite, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		accountPath(accountID)+"/invites",
		map[string]string{"email": email, "role": role},
	)
	if err != nil {
		return nil, err
	}
	var inv AccountInvite
	if err := json.Unmarshal(body, &inv); err != nil {
		return nil, fmt.Errorf("unmarshal account invite: %w", err)
	}
	return &inv, nil
}

// ListAccountInvites returns all pending invites for an account.
func (c *Client) ListAccountInvites(
	ctx context.Context, accountID string,
) ([]AccountInvite, error) {
	body, err := c.do(ctx, http.MethodGet, accountPath(accountID)+"/invites", nil)
	if err != nil {
		return nil, err
	}
	var invites []AccountInvite
	if err := json.Unmarshal(body, &invites); err != nil {
		return nil, fmt.Errorf("unmarshal account invites: %w", err)
	}
	return invites, nil
}

// RevokeAccountInvite revokes a pending account invitation.
func (c *Client) RevokeAccountInvite(ctx context.Context, accountID, inviteID string) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		accountPath(accountID)+"/invites/"+inviteID, nil,
	)
	return err
}
