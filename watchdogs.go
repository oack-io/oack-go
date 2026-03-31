package oack

import "context"

// Deprecated: Use Trigger instead.
type Watchdog = Trigger

// Deprecated: Use CreateTriggerParams instead.
type CreateWatchdogParams = CreateTriggerParams

// Deprecated: Use CreateTrigger instead.
func (c *Client) CreateWatchdog(ctx context.Context, accountID, pageID, compID string, params *CreateTriggerParams) (*Trigger, error) {
	return c.CreateTrigger(ctx, accountID, pageID, compID, params)
}

// Deprecated: Use ListTriggers instead.
func (c *Client) ListWatchdogs(ctx context.Context, accountID, pageID, compID string) ([]Trigger, error) {
	return c.ListTriggers(ctx, accountID, pageID, compID)
}

// Deprecated: Use UpdateTrigger instead.
func (c *Client) UpdateWatchdog(ctx context.Context, accountID, pageID, compID, watchdogID string, params *CreateTriggerParams) (*Trigger, error) {
	return c.UpdateTrigger(ctx, accountID, pageID, compID, watchdogID, params)
}

// Deprecated: Use DeleteTrigger instead.
func (c *Client) DeleteWatchdog(ctx context.Context, accountID, pageID, compID, watchdogID string) error {
	return c.DeleteTrigger(ctx, accountID, pageID, compID, watchdogID)
}
