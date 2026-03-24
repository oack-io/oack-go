package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents the authenticated user.
type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Provider      string `json:"provider,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Preferences holds user display preferences.
type Preferences struct {
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
	Theme      string `json:"theme"`
}

// UpdatePreferencesParams holds parameters for updating user preferences.
type UpdatePreferencesParams struct {
	Timezone   *string `json:"timezone,omitempty"`
	DateFormat *string `json:"date_format,omitempty"`
	Theme      *string `json:"theme,omitempty"`
}

// Device represents a registered push notification device.
type Device struct {
	Token     string `json:"token"`
	Platform  string `json:"platform"`
	CreatedAt string `json:"created_at"`
}

// TelegramLink holds the details for linking a Telegram account.
type TelegramLink struct {
	Code      string `json:"code"`
	BotURL    string `json:"bot_url"`
	ExpiresAt string `json:"expires_at"`
}

// TelegramLinkStatus holds the status of a Telegram link.
type TelegramLinkStatus struct {
	Linked bool   `json:"linked"`
	ChatID string `json:"chat_id"`
}

// Whoami returns the authenticated user.
func (c *Client) Whoami(ctx context.Context) (*User, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/me", nil)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}

// GetPreferences returns the authenticated user's preferences.
func (c *Client) GetPreferences(ctx context.Context) (*Preferences, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/me/preferences", nil)
	if err != nil {
		return nil, err
	}
	var p Preferences
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("unmarshal preferences: %w", err)
	}
	return &p, nil
}

// UpdatePreferences updates the authenticated user's preferences.
func (c *Client) UpdatePreferences(
	ctx context.Context, params UpdatePreferencesParams,
) (*Preferences, error) {
	body, err := c.do(ctx, http.MethodPut, "/api/v1/me/preferences", params)
	if err != nil {
		return nil, err
	}
	var p Preferences
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("unmarshal preferences: %w", err)
	}
	return &p, nil
}

// RegisterDevice registers a push notification device.
func (c *Client) RegisterDevice(
	ctx context.Context, token, platform string,
) (*Device, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/me/devices", map[string]string{
		"token":    token,
		"platform": platform,
	})
	if err != nil {
		return nil, err
	}
	var d Device
	if err := json.Unmarshal(body, &d); err != nil {
		return nil, fmt.Errorf("unmarshal device: %w", err)
	}
	return &d, nil
}

// ListDevices returns all registered push notification devices.
func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/me/devices", nil)
	if err != nil {
		return nil, err
	}
	var devices []Device
	if err := json.Unmarshal(body, &devices); err != nil {
		return nil, fmt.Errorf("unmarshal devices: %w", err)
	}
	return devices, nil
}

// UnregisterDevice removes a push notification device by token.
func (c *Client) UnregisterDevice(ctx context.Context, token string) error {
	_, err := c.do(ctx, http.MethodDelete, "/api/v1/me/devices/"+token, nil)
	return err
}

// CreateTelegramLink initiates a Telegram account linking flow.
func (c *Client) CreateTelegramLink(ctx context.Context) (*TelegramLink, error) {
	body, err := c.do(ctx, http.MethodPost, "/api/v1/me/telegram-link", nil)
	if err != nil {
		return nil, err
	}
	var tl TelegramLink
	if err := json.Unmarshal(body, &tl); err != nil {
		return nil, fmt.Errorf("unmarshal telegram link: %w", err)
	}
	return &tl, nil
}

// GetTelegramLinkStatus checks the status of a Telegram linking flow.
func (c *Client) GetTelegramLinkStatus(
	ctx context.Context, code string,
) (*TelegramLinkStatus, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/me/telegram-link/"+code+"/status", nil,
	)
	if err != nil {
		return nil, err
	}
	var s TelegramLinkStatus
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal telegram link status: %w", err)
	}
	return &s, nil
}
