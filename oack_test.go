package oack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	c := New(BearerToken("tok"))
	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL: got %q, want %q", c.baseURL, defaultBaseURL)
	}
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestNew_WithBaseURL(t *testing.T) {
	c := New(BearerToken("tok"), WithBaseURL("http://localhost:9090"))
	if c.baseURL != "http://localhost:9090" {
		t.Errorf("baseURL: got %q, want http://localhost:9090", c.baseURL)
	}
}

func TestNew_WithHTTPClient(t *testing.T) {
	hc := &http.Client{}
	c := New(BearerToken("tok"), WithHTTPClient(hc))
	if c.httpClient != hc {
		t.Error("custom httpClient not set")
	}
}

func TestBearerToken(t *testing.T) {
	tok := BearerToken("my-api-key")
	if tok.Token() != "my-api-key" {
		t.Errorf("Token(): got %q", tok.Token())
	}
}

func TestTokenFunc(t *testing.T) {
	called := false
	fn := TokenFunc(func() string {
		called = true
		return "dynamic-jwt"
	})
	if fn.Token() != "dynamic-jwt" {
		t.Errorf("Token(): got %q", fn.Token())
	}
	if !called {
		t.Error("function not called")
	}
}

func TestDo_SetsAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization: got %q, want %q", auth, "Bearer test-token")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(BearerToken("test-token"), WithBaseURL(srv.URL))
	_, err := c.do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
}

func TestDo_SetsContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type: got %q", ct)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(BearerToken("tok"), WithBaseURL(srv.URL))
	_, err := c.do(context.Background(), http.MethodPost, "/test", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("do: %v", err)
	}
}

func TestDo_NoContentTypeOnNilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "" {
			t.Errorf("Content-Type should be empty for nil body, got %q", ct)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(BearerToken("tok"), WithBaseURL(srv.URL))
	_, err := c.do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
}

func TestDo_ReturnsAPIErrorOn4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := New(BearerToken("tok"), WithBaseURL(srv.URL))
	_, err := c.do(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound, got: %v", err)
	}
}

func TestDo_ReturnsBodyOnSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"test"}`))
	}))
	defer srv.Close()

	c := New(BearerToken("tok"), WithBaseURL(srv.URL))
	body, err := c.do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	if string(body) != `{"name":"test"}` {
		t.Errorf("body: got %q", string(body))
	}
}

func TestDo_NilAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Errorf("expected no Authorization header, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(nil, WithBaseURL(srv.URL))
	_, err := c.do(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
}
