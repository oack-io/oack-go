package oack

import (
	"errors"
	"fmt"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		err  *APIError
		want string
	}{
		{&APIError{StatusCode: 404, Message: "not found"}, "oack API error (404): not found"},
		{&APIError{StatusCode: 500, Message: ""}, "oack API error (500)"},
	}
	for _, tt := range tests {
		if got := tt.err.Error(); got != tt.want {
			t.Errorf("Error(): got %q, want %q", got, tt.want)
		}
	}
}

func TestParseError_JSONError(t *testing.T) {
	err := parseError(400, []byte(`{"error":"bad input"}`))
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *APIError")
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode: got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "bad input" {
		t.Errorf("Message: got %q", apiErr.Message)
	}
}

func TestParseError_JSONMessage(t *testing.T) {
	err := parseError(422, []byte(`{"message":"invalid"}`))
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *APIError")
	}
	if apiErr.Message != "invalid" {
		t.Errorf("Message: got %q", apiErr.Message)
	}
}

func TestParseError_PlainText(t *testing.T) {
	err := parseError(500, []byte("internal error"))
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *APIError")
	}
	if apiErr.Message != "internal error" {
		t.Errorf("Message: got %q", apiErr.Message)
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(&APIError{StatusCode: 404}) {
		t.Error("expected true for 404")
	}
	if IsNotFound(&APIError{StatusCode: 403}) {
		t.Error("expected false for 403")
	}
	if IsNotFound(fmt.Errorf("other")) {
		t.Error("expected false for non-APIError")
	}
}

func TestIsForbidden(t *testing.T) {
	if !IsForbidden(&APIError{StatusCode: 403}) {
		t.Error("expected true for 403")
	}
	if IsForbidden(&APIError{StatusCode: 404}) {
		t.Error("expected false for 404")
	}
}

func TestIsConflict(t *testing.T) {
	if !IsConflict(&APIError{StatusCode: 409}) {
		t.Error("expected true for 409")
	}
}

func TestIsRateLimited(t *testing.T) {
	if !IsRateLimited(&APIError{StatusCode: 429}) {
		t.Error("expected true for 429")
	}
}

func TestIsNotFound_WrappedError(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &APIError{StatusCode: 404, Message: "gone"})
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to work through wrapping")
	}
}
