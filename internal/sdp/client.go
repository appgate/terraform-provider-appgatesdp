// Package sdp is a small, version-agnostic HTTP client for the SDP Admin API.
// It handles transport security, authentication, the versioned Accept header
// and error decoding. Request and response bodies are plain JSON; their shape
// is described per API version by the generated metadata.
package sdp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Config configures a Client.
type Config struct {
	// URL of the controller admin interface. The port defaults to 8443 and
	// the path is always /admin.
	URL          string
	Username     string
	Password     string
	Provider     string
	DeviceID     string
	BearerToken  string
	Insecure     bool
	PEMFile      string
	Timeout      time.Duration
	LoginTimeout time.Duration
	UserAgent    string
	// HTTPClient overrides the transport, mainly for tests.
	HTTPClient *http.Client
}

// Client talks to one controller using one API version.
type Client struct {
	cfg     Config
	base    *url.URL
	http    *http.Client
	version int

	mu    sync.Mutex
	token string
}

// NormalizeURL returns the admin API base URL for a controller address.
func NormalizeURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, errors.New("controller URL is required")
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("invalid controller URL %q", raw)
	}
	if u.Port() == "" {
		u.Host = net.JoinHostPort(u.Hostname(), "8443")
	}
	u.Path = "/admin"
	u.RawQuery = ""
	return u, nil
}

// New creates a client. version is the numeric API version (e.g. 25); use
// SetVersion after negotiation if it is not known yet.
func New(cfg Config, version int) (*Client, error) {
	base, err := NormalizeURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	if cfg.BearerToken == "" && (cfg.Username == "" || cfg.Password == "") {
		return nil, errors.New("username and password are required when bearer_token is not set")
	}
	if cfg.Provider == "" {
		cfg.Provider = "local"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 20 * time.Second
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		roots, _ := x509.SystemCertPool()
		if roots == nil {
			roots = x509.NewCertPool()
		}
		if cfg.PEMFile != "" {
			pem, err := os.ReadFile(cfg.PEMFile)
			if err != nil {
				return nil, fmt.Errorf("read pem_filepath: %w", err)
			}
			if !roots.AppendCertsFromPEM(pem) {
				return nil, fmt.Errorf("pem_filepath %s contains no certificates", cfg.PEMFile)
			}
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{RootCAs: roots, InsecureSkipVerify: cfg.Insecure} //nolint:gosec // insecure is an explicit, documented provider option.
		transport.TLSHandshakeTimeout = cfg.Timeout
		httpClient = &http.Client{Transport: transport, Timeout: 2 * cfg.Timeout}
	}
	c := &Client{cfg: cfg, base: base, http: httpClient, version: version}
	if cfg.BearerToken != "" {
		c.token = cfg.BearerToken
	}
	return c, nil
}

// SetVersion sets the API version used for subsequent requests.
func (c *Client) SetVersion(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.version = v
}

func (c *Client) acceptHeader() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return fmt.Sprintf("application/vnd.appgate.peer-v%d+json", c.version)
}

// Error is an error response of the Admin API.
type Error struct {
	Status  int          `json:"-"`
	Method  string       `json:"-"`
	Path    string       `json:"-"`
	ID      string       `json:"id"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
	// Min and max supported API versions, returned with 406 responses.
	MinSupportedVersion int `json:"minSupportedVersion"`
	MaxSupportedVersion int `json:"maxSupportedVersion"`
}

// FieldError is one entry of a 422 validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s: HTTP %d", e.Method, e.Path, e.Status)
	if e.ID != "" {
		fmt.Fprintf(&b, " %s", e.ID)
	}
	if e.Message != "" {
		fmt.Fprintf(&b, ": %s", e.Message)
	}
	for _, fe := range e.Errors {
		fmt.Fprintf(&b, "\n  %s: %s", fe.Field, fe.Message)
	}
	return b.String()
}

// IsNotFound reports whether err is an HTTP 404 from the Admin API.
func IsNotFound(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Status == http.StatusNotFound
}

// Do sends a request and decodes the JSON response into out (if non-nil).
// The client logs in on first use and once more if the token was rejected.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	token, err := c.ensureToken(ctx)
	if err != nil {
		return err
	}
	accept := c.acceptHeader()
	err = c.do(ctx, accept, method, path, query, body, out, token)
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusUnauthorized && c.cfg.BearerToken == "" {
		c.mu.Lock()
		if c.token == token {
			c.token = ""
		}
		c.mu.Unlock()
		if token, err = c.ensureToken(ctx); err != nil {
			return err
		}
		return c.do(ctx, accept, method, path, query, body, out, token)
	}
	return err
}

func (c *Client) do(ctx context.Context, accept, method, path string, query url.Values, body, out any, token string) error {
	// path is already escaped by the caller (object ids may contain ',' or
	// '='); keep it as RawPath so it is not escaped a second time.
	u := *c.base
	u.RawPath = strings.TrimSuffix(c.base.Path, "/") + path
	unescaped, err := url.PathUnescape(u.RawPath)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	u.Path = unescaped
	u.RawQuery = query.Encode()
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", accept)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.cfg.UserAgent != "" {
		req.Header.Set("User-Agent", c.cfg.UserAgent)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		var unknownCA x509.UnknownAuthorityError
		if errors.As(err, &unknownCA) {
			return fmt.Errorf("%s %s: %w; set pem_filepath to the controller CA or enable insecure", method, path, err)
		}
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s %s: read response: %w", method, path, err)
	}
	if resp.StatusCode >= 300 {
		apiErr := &Error{Status: resp.StatusCode, Method: method, Path: path}
		if json.Unmarshal(data, apiErr) != nil && len(data) > 0 {
			apiErr.Message = strings.TrimSpace(string(data))
		}
		return apiErr
	}
	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("%s %s: decode response: %w", method, path, err)
	}
	return nil
}

func (c *Client) ensureToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" {
		return c.token, nil
	}
	token, err := c.login(ctx)
	if err != nil {
		return "", err
	}
	c.token = token
	return token, nil
}

type loginRequest struct {
	ProviderName string `json:"providerName"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	DeviceID     string `json:"deviceId"`
}

// login authenticates with username and password. The controller may still
// be starting when Terraform runs (e.g. created in the same apply), so
// connection errors and 5xx responses are retried until LoginTimeout.
// Must be called with c.mu held.
func (c *Client) login(ctx context.Context) (string, error) {
	req := loginRequest{ProviderName: c.cfg.Provider, Username: c.cfg.Username, Password: c.cfg.Password, DeviceID: c.cfg.DeviceID}
	deadline := time.Now().Add(c.cfg.LoginTimeout)
	delay := 500 * time.Millisecond
	accept := fmt.Sprintf("application/vnd.appgate.peer-v%d+json", c.version)
	for {
		var resp struct {
			Token string `json:"token"`
		}
		err := c.do(ctx, accept, http.MethodPost, "/login", nil, req, &resp, "")
		if err == nil {
			if resp.Token == "" {
				return "", errors.New("POST /login: response contains no token")
			}
			return resp.Token, nil
		}
		var apiErr *Error
		retryable := !errors.As(err, &apiErr) || apiErr.Status >= 500
		var unknownCA x509.UnknownAuthorityError
		if errors.As(err, &unknownCA) {
			retryable = false
		}
		if !retryable || time.Now().Add(delay).After(deadline) {
			return "", fmt.Errorf("login failed: %w", err)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
		}
		delay = min(delay*3/2, 30*time.Second)
	}
}

// SupportedVersions asks the controller which API versions it accepts. The
// login endpoint answers an unsupported Accept header with HTTP 406 and the
// supported range.
func (c *Client) SupportedVersions(ctx context.Context) (minVersion, maxVersion int, err error) {
	req := loginRequest{ProviderName: c.cfg.Provider, DeviceID: c.cfg.DeviceID}
	err = c.do(ctx, "application/vnd.appgate.peer-v5+json", http.MethodPost, "/login", nil, req, nil, "")
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotAcceptable && apiErr.MaxSupportedVersion > 0 {
		return apiErr.MinSupportedVersion, apiErr.MaxSupportedVersion, nil
	}
	if err == nil {
		err = errors.New("controller accepted an unsupported API version")
	}
	return 0, 0, fmt.Errorf("determine supported API versions: %w", err)
}
