// Package oack provides a Go client for the Oack monitoring API.
//
// Create a client with a static token (API key):
//
//	c := oack.New(oack.BearerToken("sk-..."))
//
// Or with a dynamic token (refreshable JWT):
//
//	c := oack.New(oack.TokenFunc(myTokenProvider))
package oack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.oack.io"

// AuthMethod provides the Authorization bearer token.
type AuthMethod interface {
	Token() string
}

// BearerToken authenticates with a static token (API key or JWT).
type BearerToken string

func (t BearerToken) Token() string { return string(t) }

// TokenFunc authenticates with a dynamic token (e.g. refreshable JWT).
type TokenFunc func() string

func (f TokenFunc) Token() string { return f() }

// Client is the Oack API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       AuthMethod
	userAgent  string
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL (https://api.oack.io).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient sets a custom http.Client for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets the User-Agent header for all requests.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// New creates an Oack API client.
func New(auth AuthMethod, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		auth: auth,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// do executes an HTTP request and returns the response body.
// Non-2xx responses are returned as *APIError.
func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if c.auth != nil {
		if tok := c.auth.Token(); tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
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

	return respBody, nil
}
