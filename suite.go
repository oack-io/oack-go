package oack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// DeploySuiteOpts holds optional filters for a suite deployment.
type DeploySuiteOpts struct {
	PwProject string
	PwGrep    string
	PwTag     string
	GitSHA     string
	GitBranch  string
	GitOrigin  string
	DeployHost string
	DeployCmd  string
}

// DeploySuiteResult holds the response from a suite upload.
type DeploySuiteResult struct {
	SuiteURL  string   `json:"suite_url"`
	DepsURL   string   `json:"deps_url"`
	SizeBytes int      `json:"size_bytes"`
	TestFiles []string `json:"test_files"`
}

// DeploySuite uploads a Playwright project tarball to a browser monitor.
func (c *Client) DeploySuite(
	ctx context.Context, teamID, monitorID string,
	tarball []byte, opts *DeploySuiteOpts,
) (*DeploySuiteResult, error) {
	path := monitorPath(teamID, monitorID) + "/suite"

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Suite file.
	fw, err := w.CreateFormFile("suite", "suite.tar.gz")
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := fw.Write(tarball); err != nil {
		return nil, fmt.Errorf("write tarball: %w", err)
	}

	// Optional filters and metadata.
	if opts != nil {
		if opts.PwProject != "" {
			_ = w.WriteField("pw_project", opts.PwProject)
		}
		if opts.PwGrep != "" {
			_ = w.WriteField("pw_grep", opts.PwGrep)
		}
		if opts.PwTag != "" {
			_ = w.WriteField("pw_tag", opts.PwTag)
		}
		if opts.GitSHA != "" {
			_ = w.WriteField("git_sha", opts.GitSHA)
		}
		if opts.GitBranch != "" {
			_ = w.WriteField("git_branch", opts.GitBranch)
		}
		if opts.GitOrigin != "" {
			_ = w.WriteField("git_origin", opts.GitOrigin)
		}
		if opts.DeployHost != "" {
			_ = w.WriteField("deploy_host", opts.DeployHost)
		}
		if opts.DeployCmd != "" {
			_ = w.WriteField("deploy_cmd", opts.DeployCmd)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if c.auth != nil {
		if tok := c.auth.Token(); tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, respBody)
	}

	var result DeploySuiteResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &result, nil
}
