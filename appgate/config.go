package appgate

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"
	"github.com/hashicorp/go-version"

	"github.com/google/uuid"
)

const (
	// DefaultDescription is the default string for terraform resources.
	DefaultDescription = "Managed by terraform"
)

// Config for appgate provider.
type Config struct {
	URL      string
	Username string
	Password string
	Provider string
	Insecure bool
	Timeout  int
	Debug    bool
	Version  int
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

// AppgateConfigFile used with  "config_path"
type AppgateConfigFile struct {
	URL           string `json:"appgate_url"`
	Username      string `json:"appgate_username"`
	Password      string `json:"appgate_password"`
	Provider      string `json:"appgate_provider"`
	ClientVersion int    `json:"appgate_client_version"`
	Insecure      bool   `json:"appgate_insecure"`
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
	token := fmt.Sprintf("Bearer %s", *openapi.PtrString(*response.Token))

	return token, nil
}

func (c *Client) login() (*openapi.LoginResponse, error) {
	ctx := context.Background()
	loginOpts := openapi.LoginRequest{
		ProviderName: c.Config.Provider,
		Username:     openapi.PtrString(c.Config.Username),
		Password:     openapi.PtrString(c.Config.Password),
		DeviceId:     uuid.New().String(),
	}

	loginResponse, _, err := c.API.LoginApi.LoginPost(ctx).LoginRequest(loginOpts).Execute()
	if err != nil {
		return nil, prettyPrintAPIError(err)
	}
	return &loginResponse, nil
}
