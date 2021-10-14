package appgate

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-version"

	"github.com/google/uuid"
)

const (
	// DefaultDescription is the default string for terraform resources.
	DefaultDescription = "Managed by terraform"
)

// Config for appgate provider.
type Config struct {
	URL          string `json:"appgate_url,omitempty"`
	Username     string `json:"appgate_username,omitempty"`
	Password     string `json:"appgate_password,omitempty"`
	Provider     string `json:"appgate_provider,omitempty"`
	Insecure     bool   `json:"appgate_insecure,omitempty"`
	Timeout      int    `json:"appgate_timeout,omitempty"`
	LoginTimeout int    `json:"appgate_login_timeout,omitempty"`
	Debug        bool   `json:"appgate_http_debug,omitempty"`
	Version      int    `json:"appgate_client_version,omitempty"`
}

// Client is the appgate API client.
type Client struct {
	Token                  string
	UUID                   string
	ApplianceVersion       *version.Version
	LatestSupportedVersion *version.Version
	ClientVersion          int
	API                    *openapi.APIClient
	Config                 *Config
}

// Client creates
func (c *Config) Client() (*Client, error) {
	timeoutDuration := time.Duration(c.Timeout)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.Insecure,
		},
		Dial: (&net.Dialer{
			Timeout: timeoutDuration * time.Second,
		}).Dial,
		TLSHandshakeTimeout: timeoutDuration * time.Second,
	}

	httpclient := &http.Client{
		Transport: tr,
		Timeout:   ((timeoutDuration * 2) * time.Second),
	}
	clientCfg := &openapi.Configuration{
		DefaultHeader: map[string]string{
			"Accept": fmt.Sprintf("application/vnd.appgate.peer-v%d+json", c.Version),
		},
		UserAgent: "Appgate-TerraformProvider/1.0.0/go",
		Debug:     c.Debug,
		Servers: []openapi.ServerConfiguration{
			{
				URL: c.URL,
			},
		},
		HTTPClient: httpclient,
	}
	apiClient := openapi.NewAPIClient(clientCfg)
	if apiClient == nil {
		return nil, errors.New("failed to initialize api client")
	}

	client := &Client{
		API:           apiClient,
		ClientVersion: c.Version,
		Config:        c,
	}

	return client, nil
}

func guessVersion(response *openapi.LoginResponse, clientVersion int) (*version.Version, error) {
	if response.HasVersion() {
		currentVersion, err := version.NewVersion(*response.Version)
		if err != nil {
			return nil, err
		}
		return currentVersion, nil
	}
	switch clientVersion {
	case Version15:
		return version.NewVersion("5.4.0+estimated")
	case Version16:
		return version.NewVersion("5.5.0+estimated")

	}
	return nil, fmt.Errorf("could not determine appliance version with client version %d", clientVersion)
}

// GetToken makes first login and initiate the client towards the controller.
// this is always the first made
func (c *Client) GetToken() (string, error) {
	if len(c.Token) > 0 {
		return c.Token, nil
	}
	cfg := c.Config
	response, err := c.login()
	if err != nil {
		return "", err
	}

	latestSupportedVersion, err := version.NewVersion(ApplianceVersionMap[DefaultClientVersion])
	if err != nil {
		return "", err
	}

	currentVersion, err := guessVersion(response, cfg.Version)
	if err != nil {
		return "", err
	}
	c.ApplianceVersion = currentVersion
	c.LatestSupportedVersion = latestSupportedVersion
	c.Token = fmt.Sprintf("Bearer %s", *openapi.PtrString(*response.Token))

	return c.Token, nil
}

var exponentialBackOff = backoff.ExponentialBackOff{
	InitialInterval:     500 * time.Millisecond,
	RandomizationFactor: 0.5,
	Multiplier:          1.5,
	MaxInterval:         30 * time.Second,
	Stop:                backoff.Stop,
	Clock:               backoff.SystemClock,
}

func (c *Client) loginTimoutDuration() time.Duration {
	// This is just intend to be used within unit tests.
	if c.Config.LoginTimeout > 0 {
		return time.Duration(c.Config.LoginTimeout) * time.Second
	}
	return 5 * time.Minute
}

func (c *Client) login() (*openapi.LoginResponse, error) {
	ctx := context.Background()
	loginOpts := openapi.LoginRequest{
		ProviderName: c.Config.Provider,
		Username:     openapi.PtrString(c.Config.Username),
		Password:     openapi.PtrString(c.Config.Password),
		DeviceId:     uuid.New().String(),
	}

	// Since /login is the first request we do, it provide us the earliest check if a controller is up and running
	// if the appgatesdp provider is combined with, for example  aws_instance where we create the inital controller
	// it might take awhile for the controller to startup and be responsive, so until its up it can return 500, 502, 503
	// these status code is treated as retryable errors, during exponentialBackOff.MaxElapsedTime window.
	// we will use this exponential backoff to retry until we get a 200-400 HTTP response from /login
	exponentialBackOff.MaxElapsedTime = c.loginTimoutDuration()
	loginResponse := &openapi.LoginResponse{}
	err := backoff.Retry(func() error {
		login, response, err := c.API.LoginApi.LoginPost(ctx).LoginRequest(loginOpts).Execute()
		if response == nil {
			log.Printf("[DEBUG] Login failed, No response %s", err)
			return fmt.Errorf("No response from controller %w", err)
		}
		if response.StatusCode >= 500 {
			log.Printf("[DEBUG] Login failed, controller not responding, got HTTP %d", response.StatusCode)
			return fmt.Errorf("Controller got %w", err)
		}
		if err != nil {
			log.Printf("[DEBUG] Login failed permanently, got HTTP %d", response.StatusCode)
			return &backoff.PermanentError{Err: err}
		}
		loginResponse = &login
		return nil
	}, &exponentialBackOff)
	if err != nil {
		return nil, prettyPrintAPIError(err)
	}
	log.Printf("[DEBUG] Login OK")
	return loginResponse, nil
}
