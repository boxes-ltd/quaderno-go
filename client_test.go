package quaderno

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_defaults(t *testing.T) {
	c := NewClient("key", "https://api.example.com")

	if c.apiKey != "key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "key")
	}
	if c.apiUrl != "https://api.example.com" {
		t.Errorf("apiUrl = %q, want %q", c.apiUrl, "https://api.example.com")
	}
	if c.apiVersion != DefaultAPIVersion {
		t.Errorf("apiVersion = %q, want %q", c.apiVersion, DefaultAPIVersion)
	}
	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
	if c.Taxes == nil {
		t.Error("Taxes service should not be nil")
	}
	if c.Transactions == nil {
		t.Error("Transactions service should not be nil")
	}
}

func TestNewClient_stripsTrailingSlash(t *testing.T) {
	c := NewClient("key", "https://api.example.com/")
	if c.apiUrl != "https://api.example.com" {
		t.Errorf("apiUrl = %q, want trailing slash stripped", c.apiUrl)
	}
}

func TestNewClient_withApiVersion(t *testing.T) {
	c := NewClient("key", "https://api.example.com", WithApiVersion("20230101"))
	if c.apiVersion != "20230101" {
		t.Errorf("apiVersion = %q, want %q", c.apiVersion, "20230101")
	}
}

func TestNewClient_withHttpClient(t *testing.T) {
	custom := &http.Client{}
	c := NewClient("key", "https://api.example.com", WithHttpClient(custom))
	// When logLevel is None the custom client is used directly.
	if c.httpClient != custom {
		t.Error("expected custom http client to be used")
	}
}

func TestNewClient_withUserAgent(t *testing.T) {
	c := NewClient("key", "https://api.example.com", WithUserAgent("my-agent/1.0"))
	if c.userAgent != "my-agent/1.0" {
		t.Errorf("userAgent = %q, want %q", c.userAgent, "my-agent/1.0")
	}
}

func TestNewClient_withLogLevel_wrapsTransport(t *testing.T) {
	c := NewClient("key", "https://api.example.com", WithLogLevel(LogLevelBasic))
	transport, ok := c.httpClient.Transport.(*httpLogger)
	if !ok {
		t.Fatal("expected http.Client.Transport to be *httpLogger")
	}
	if transport.level != LogLevelBasic {
		t.Errorf("httpLogger level = %v, want %v", transport.level, LogLevelBasic)
	}
}

func TestDoRequest_setsAuthAndHeaders(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient("mykey", srv.URL, WithApiVersion("20991231"), WithUserAgent("test-agent"))
	if err := c.doRequest(context.Background(), http.MethodGet, "/ping", nil, nil, nil); err != nil {
		t.Fatalf("doRequest() unexpected error: %v", err)
	}

	user, pass, ok := capturedReq.BasicAuth()
	if !ok {
		t.Fatal("expected basic auth to be set")
	}
	if user != "mykey" {
		t.Errorf("basic auth user = %q, want %q", user, "mykey")
	}
	if pass != "x" {
		t.Errorf("basic auth pass = %q, want %q", pass, "x")
	}

	wantAccept := "application/json; api_version: 20991231"
	if got := capturedReq.Header.Get("Accept"); got != wantAccept {
		t.Errorf("Accept = %q, want %q", got, wantAccept)
	}
	if got := capturedReq.Header.Get("User-Agent"); got != "test-agent" {
		t.Errorf("User-Agent = %q, want %q", got, "test-agent")
	}
}

func TestDoRequest_setsContentTypeForBody(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	if err := c.doRequest(context.Background(), http.MethodPost, "/test", nil, map[string]string{"k": "v"}, nil); err != nil {
		t.Fatalf("doRequest() unexpected error: %v", err)
	}

	if got := capturedReq.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestDoRequest_decodesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"foo": "bar"})
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	var out map[string]string
	if err := c.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foo"] != "bar" {
		t.Errorf("decoded foo = %q, want %q", out["foo"], "bar")
	}
}

func TestDoRequest_returnsApiErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"invalid"}`))
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	err := c.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, nil)

	apiErr, ok := errors.AsType[*ApiError](err)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnprocessableEntity)
	}
	if string(apiErr.Body) != `{"error":"invalid"}` {
		t.Errorf("Body = %q, want %q", apiErr.Body, `{"error":"invalid"}`)
	}
}

func TestApiError_Error(t *testing.T) {
	err := &ApiError{StatusCode: 422, Body: []byte("bad input")}
	want := "API error (status 422): bad input"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestPing_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			t.Errorf("path = %q, want /ping", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	if err := c.Ping(context.Background()); err != nil {
		t.Errorf("Ping() unexpected error: %v", err)
	}
}

func TestPing_failure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := NewClient("key", srv.URL)
	if err := c.Ping(context.Background()); err == nil {
		t.Error("Ping() expected error, got nil")
	}
}