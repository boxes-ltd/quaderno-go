package quaderno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultAPIVersion = "20260309"

func defaultUserAgent() string {
	userAgent := "quaderno-go"
	version := getVersion()
	if version != "" {
		userAgent += "/" + version
	}
	return userAgent
}

type Client struct {
	apiKey     string
	apiUrl     string
	apiVersion string
	httpClient *http.Client
	logLevel   LogLevel
	userAgent  string
}

type Option func(*Client)

func WithApiVersion(apiVersion string) Option {
	return func(c *Client) {
		c.apiVersion = apiVersion
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithLogLevel(level LogLevel) Option {
	return func(c *Client) {
		c.logLevel = level
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func NewClient(apiKey, apiUrl string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		apiUrl:     strings.TrimSuffix(apiUrl, "/"),
		apiVersion: DefaultAPIVersion,
		httpClient: http.DefaultClient,
		userAgent:  defaultUserAgent(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.logLevel != LogLevelNone {
		clonedClient := *c.httpClient
		baseTransport := clonedClient.Transport
		if baseTransport == nil {
			baseTransport = http.DefaultTransport
		}
		clonedClient.Transport = &httpLogger{
			transport: baseTransport,
			level:     c.logLevel,
		}
		c.httpClient = &clonedClient
	}
	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body, response any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	reqUrl, err := url.Parse(c.apiUrl + path)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	if query != nil {
		reqUrl.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		reqBody = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "x")
	req.Header.Set("Accept", "application/json; api_version: "+c.apiVersion)
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	err := c.doRequest(ctx, http.MethodGet, "/ping", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
