package oack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// APIError represents a non-2xx response from the Oack API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("oack API error (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("oack API error (%d)", e.StatusCode)
}

func parseError(statusCode int, body []byte) error {
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil {
		msg := errResp.Error
		if msg == "" {
			msg = errResp.Message
		}
		if msg != "" {
			return &APIError{StatusCode: statusCode, Message: msg}
		}
	}
	return &APIError{StatusCode: statusCode, Message: string(body)}
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool { return hasStatus(err, http.StatusNotFound) }

// IsForbidden returns true if the error is a 403 Forbidden.
func IsForbidden(err error) bool { return hasStatus(err, http.StatusForbidden) }

// IsConflict returns true if the error is a 409 Conflict.
func IsConflict(err error) bool { return hasStatus(err, http.StatusConflict) }

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func IsRateLimited(err error) bool { return hasStatus(err, http.StatusTooManyRequests) }

func hasStatus(err error, code int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == code
}
