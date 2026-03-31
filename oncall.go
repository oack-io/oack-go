package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OnCallSchedule represents an on-call rotation schedule.
type OnCallSchedule struct {
	ID           string   `json:"id"`
	AccountID    string   `json:"account_id"`
	Name         string   `json:"name"`
	Timezone     string   `json:"timezone"`
	RotationType string   `json:"rotation_type"`
	Participants []string `json:"participants"`
	HandoffTime  string   `json:"handoff_time"`
	HandoffDay   int      `json:"handoff_day"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// CreateScheduleParams holds parameters for creating an on-call schedule.
type CreateScheduleParams struct {
	Name         string   `json:"name"`
	Timezone     string   `json:"timezone,omitempty"`
	RotationType string   `json:"rotation_type,omitempty"`
	Participants []string `json:"participants,omitempty"`
	HandoffTime  string   `json:"handoff_time,omitempty"`
	HandoffDay   int      `json:"handoff_day,omitempty"`
}

// UpdateScheduleParams holds parameters for updating an on-call schedule.
type UpdateScheduleParams struct {
	Name         string   `json:"name,omitempty"`
	Timezone     string   `json:"timezone,omitempty"`
	RotationType string   `json:"rotation_type,omitempty"`
	Participants []string `json:"participants,omitempty"`
	HandoffTime  string   `json:"handoff_time,omitempty"`
	HandoffDay   int      `json:"handoff_day,omitempty"`
}

// OnCallOverride represents a temporary on-call schedule override.
type OnCallOverride struct {
	ID                string `json:"id"`
	ScheduleID        string `json:"schedule_id"`
	OriginalUserID    string `json:"original_user_id"`
	ReplacementUserID string `json:"replacement_user_id"`
	StartAt           string `json:"start_at"`
	EndAt             string `json:"end_at"`
	Reason            string `json:"reason"`
	CreatedAt         string `json:"created_at"`
}

// CreateOverrideParams holds parameters for creating an on-call override.
type CreateOverrideParams struct {
	OriginalUserID    string `json:"original_user_id"`
	ReplacementUserID string `json:"replacement_user_id"`
	StartAt           string `json:"start_at"`
	EndAt             string `json:"end_at"`
	Reason            string `json:"reason,omitempty"`
}

// WhosOnCall represents who is currently on call for a schedule.
type WhosOnCall struct {
	ScheduleID   string `json:"schedule_id"`
	ScheduleName string `json:"schedule_name"`
	UserID       string `json:"user_id"`
	OverrideID   string `json:"override_id,omitempty"`
}

// scheduleBasePath returns the base URL path for on-call schedule operations.
func scheduleBasePath(accountID string) string {
	return "/api/v1/accounts/" + accountID + "/oncall/schedules"
}

// schedulePath returns the URL path for a specific on-call schedule.
func schedulePath(accountID, scheduleID string) string {
	return scheduleBasePath(accountID) + "/" + scheduleID
}

// ── Schedules

// CreateSchedule creates a new on-call schedule.
func (c *Client) CreateSchedule(
	ctx context.Context, accountID string, params *CreateScheduleParams,
) (*OnCallSchedule, error) {
	body, err := c.do(ctx, http.MethodPost, scheduleBasePath(accountID), params)
	if err != nil {
		return nil, err
	}
	var s OnCallSchedule
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal schedule: %w", err)
	}
	return &s, nil
}

// GetSchedule returns a single on-call schedule by ID.
func (c *Client) GetSchedule(
	ctx context.Context, accountID, scheduleID string,
) (*OnCallSchedule, error) {
	body, err := c.do(ctx, http.MethodGet, schedulePath(accountID, scheduleID), nil)
	if err != nil {
		return nil, err
	}
	var s OnCallSchedule
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal schedule: %w", err)
	}
	return &s, nil
}

// ListSchedules returns all on-call schedules for an account.
func (c *Client) ListSchedules(
	ctx context.Context, accountID string,
) ([]OnCallSchedule, error) {
	body, err := c.do(ctx, http.MethodGet, scheduleBasePath(accountID), nil)
	if err != nil {
		return nil, err
	}
	var schedules []OnCallSchedule
	if err := json.Unmarshal(body, &schedules); err != nil {
		return nil, fmt.Errorf("unmarshal schedules: %w", err)
	}
	return schedules, nil
}

// UpdateSchedule updates an existing on-call schedule.
func (c *Client) UpdateSchedule(
	ctx context.Context, accountID, scheduleID string, params *UpdateScheduleParams,
) (*OnCallSchedule, error) {
	body, err := c.do(ctx, http.MethodPut, schedulePath(accountID, scheduleID), params)
	if err != nil {
		return nil, err
	}
	var s OnCallSchedule
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("unmarshal schedule: %w", err)
	}
	return &s, nil
}

// DeleteSchedule deletes an on-call schedule by ID.
func (c *Client) DeleteSchedule(ctx context.Context, accountID, scheduleID string) error {
	_, err := c.do(ctx, http.MethodDelete, schedulePath(accountID, scheduleID), nil)
	return err
}

// ── Overrides

// CreateOverride creates a temporary on-call override on a schedule.
func (c *Client) CreateOverride(
	ctx context.Context, accountID, scheduleID string, params *CreateOverrideParams,
) (*OnCallOverride, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		schedulePath(accountID, scheduleID)+"/overrides", params,
	)
	if err != nil {
		return nil, err
	}
	var o OnCallOverride
	if err := json.Unmarshal(body, &o); err != nil {
		return nil, fmt.Errorf("unmarshal override: %w", err)
	}
	return &o, nil
}

// ListOverrides returns all overrides for a schedule.
func (c *Client) ListOverrides(
	ctx context.Context, accountID, scheduleID string,
) ([]OnCallOverride, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		schedulePath(accountID, scheduleID)+"/overrides", nil,
	)
	if err != nil {
		return nil, err
	}
	var overrides []OnCallOverride
	if err := json.Unmarshal(body, &overrides); err != nil {
		return nil, fmt.Errorf("unmarshal overrides: %w", err)
	}
	return overrides, nil
}

// DeleteOverride deletes an on-call override by ID.
func (c *Client) DeleteOverride(
	ctx context.Context, accountID, scheduleID, overrideID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		schedulePath(accountID, scheduleID)+"/overrides/"+overrideID, nil,
	)
	return err
}

// ── Who's on call

// GetWhosOnCall returns the currently on-call user for each schedule in the account.
func (c *Client) GetWhosOnCall(
	ctx context.Context, accountID string,
) ([]WhosOnCall, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/accounts/"+accountID+"/oncall/now", nil,
	)
	if err != nil {
		return nil, err
	}
	var info []WhosOnCall
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("unmarshal whos on call: %w", err)
	}
	return info, nil
}
