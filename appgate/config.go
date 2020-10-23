package appgate

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/google/uuid"
)

const (
	// Version is the Appgate controller version.
	Version = 13

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
}

// Client is the appgate API client.
type Client struct {
	Token string
	UUID  string
	API   *openapi.APIClient
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
		// Host:   c.URL,
		// Scheme: "https",
		DefaultHeader: map[string]string{
			"Accept": fmt.Sprintf("application/vnd.appgate.peer-v%d+json", Version),
		},
		UserAgent: "Appgate-TerraformProvider/1.0.0/go",
		Debug:     c.Debug,
		Servers: []openapi.ServerConfiguration{
			{
				URL:         c.URL,
				Description: "Controller one",
			},
		},
		HTTPClient: httpclient,
	}
	apiClient := openapi.NewAPIClient(clientCfg)

	token, err := getToken(apiClient, c)
	if err != nil {
		return nil, err
	}

	client := &Client{
		API:   apiClient,
		Token: token,
	}
	return client, nil
}

func getToken(apiClient *openapi.APIClient, cfg *Config) (string, error) {

	ctx := context.Background()
	// Login first, save token
	loginOpts := openapi.LoginRequest{
		ProviderName: cfg.Provider,
		Username:     openapi.PtrString(cfg.Username),
		Password:     openapi.PtrString(cfg.Password),
		DeviceId:     uuid.New().String(),
	}

	loginResponse, _, err := apiClient.LoginApi.LoginPost(ctx).LoginRequest(loginOpts).Execute()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", *openapi.PtrString(*loginResponse.Token)), nil
}
