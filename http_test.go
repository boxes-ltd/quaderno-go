package quaderno

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.fn(r)
}

func okResponse() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
	}
}

func TestHttpLogger_masksAuthorizationHeader(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	t.Cleanup(func() { log.SetOutput(io.Discard) })

	logger := &httpLogger{
		level: LogLevelHeaders,
		transport: &mockRoundTripper{fn: func(r *http.Request) (*http.Response, error) {
			return okResponse(), nil
		}},
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/ping", nil)
	req.SetBasicAuth("myapikey", "x")

	_, err := logger.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}

	output := logOutput.String()
	if strings.Contains(output, "myapikey") {
		t.Error("log output should not contain the raw API key")
	}
	if !strings.Contains(output, "Basic ****") {
		t.Errorf("log output should contain masked Basic token, got: %s", output)
	}
}

func TestHttpLogger_logLevelNone_skipsLogging(t *testing.T) {
	called := false
	logger := &httpLogger{
		level: LogLevelNone,
		transport: &mockRoundTripper{fn: func(r *http.Request) (*http.Response, error) {
			called = true
			return okResponse(), nil
		}},
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/ping", nil)
	_, err := logger.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	if !called {
		t.Error("expected underlying transport to be called")
	}
}

func TestHttpLogger_body_logsRequestBody(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	t.Cleanup(func() { log.SetOutput(io.Discard) })

	logger := &httpLogger{
		level: LogLevelBody,
		transport: &mockRoundTripper{fn: func(r *http.Request) (*http.Response, error) {
			return okResponse(), nil
		}},
	}

	body := `{"hello":"world"}`
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/test", strings.NewReader(body))
	_, err := logger.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}

	if !strings.Contains(logOutput.String(), `{"hello":"world"}`) {
		t.Errorf("expected request body in log output, got: %s", logOutput.String())
	}
}