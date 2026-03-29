package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TestScriptParams holds parameters for a one-off script test run.
type TestScriptParams struct {
	Script       string            `json:"script,omitempty"`
	Suite        string            `json:"suite,omitempty"` // base64 tar.gz for test_suite mode
	PwProject    string            `json:"pw_project,omitempty"`
	PwGrep       string            `json:"pw_grep,omitempty"`
	EnvOverrides map[string]string `json:"env_overrides,omitempty"`
}

// TestScriptResult holds the result of a one-off script test run.
type TestScriptResult struct {
	Passed          bool             `json:"passed"`
	TotalMs         int64            `json:"total_ms"`
	Error           string           `json:"error,omitempty"`
	Status          int              `json:"status"`
	Steps           []StepResult     `json:"steps,omitempty"`
	ConsoleMessages []ConsoleMessage `json:"console_messages,omitempty"`
	ScreenshotURL   string           `json:"screenshot_url,omitempty"`
	ReportURL       string           `json:"report_url,omitempty"`
	TestCount       int              `json:"test_count,omitempty"`
	PassCount       int              `json:"pass_count,omitempty"`
	FailCount       int              `json:"fail_count,omitempty"`
	SkipCount       int              `json:"skip_count,omitempty"`
	WebVitals       *WebVitals       `json:"web_vitals,omitempty"`
}

// WebVitals holds core web vitals metrics.
type WebVitals struct {
	LcpMs  float64 `json:"lcp_ms"`
	FcpMs  float64 `json:"fcp_ms"`
	Cls    float64 `json:"cls"`
	TtfbMs float64 `json:"ttfb_ms"`
}

type testScriptSubmitResponse struct {
	TestID string `json:"test_id"`
}

type testScriptPollResponse struct {
	Status string            `json:"status"`
	Result *TestScriptResult `json:"result,omitempty"`
}

// TestScript submits a one-off browser script test and polls for the result.
// The script is sent to the server, which pushes it to a browser-capable
// checker. This method polls until the result is available or the context
// is cancelled (up to 5 minutes).
func (c *Client) TestScript(
	ctx context.Context, teamID, monitorID string, params *TestScriptParams,
) (*TestScriptResult, error) {
	// 1. Submit the test.
	submitPath := monitorPath(teamID, monitorID) + "/test-script"
	body, err := c.do(ctx, http.MethodPost, submitPath, params)
	if err != nil {
		return nil, err
	}
	var submit testScriptSubmitResponse
	if err := json.Unmarshal(body, &submit); err != nil {
		return nil, fmt.Errorf("unmarshal submit response: %w", err)
	}
	if submit.TestID == "" {
		return nil, fmt.Errorf("server returned empty test_id")
	}

	// 2. Poll for result.
	pollPath := submitPath + "/" + submit.TestID
	pollInterval := 2 * time.Second
	timeout := 5 * time.Minute

	deadline := time.After(timeout)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline:
			return nil, fmt.Errorf("test timed out after %s", timeout)
		case <-ticker.C:
			body, err := c.do(ctx, http.MethodGet, pollPath, nil)
			if err != nil {
				return nil, err
			}
			var poll testScriptPollResponse
			if err := json.Unmarshal(body, &poll); err != nil {
				return nil, fmt.Errorf("unmarshal poll response: %w", err)
			}
			if poll.Status == "done" && poll.Result != nil {
				return poll.Result, nil
			}
			// Still running — continue polling.
		}
	}
}
