package appgate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
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
	BearerToken  string `json:"appgate_bearer_token,omitempty"`
	PemFilePath  string `json:"appgate_pem_filepath,omitempty"`
}

// Validate makes sure we have minimum required configuration values to authenticate against the controller.
func (c *Config) Validate(usingFile bool) error {
	// we won't validate the configuration if we are using the config_file
	// this is because we want defer it until the file has been populated.
	if usingFile {
		return nil
	}
	if !isUrl(c.URL) {
		return fmt.Errorf("Controller URL is mandatory, got %q", c.URL)
	}
	if len(c.BearerToken) > 0 {
		_, err := b64.StdEncoding.DecodeString(c.BearerToken)
		if err != nil {
			return fmt.Errorf("appgate bearer_token set, but invalid format, expected base64 %s", err)
		}
	} else if len(c.Username) < 1 && len(c.Password) < 1 {
		return fmt.Errorf("username and password required if appgate bearer token is empty")
	}
	keys := make([]int, 0, len(ApplianceVersionMap))
	for k := range ApplianceVersionMap {
		keys = append(keys, k)
	}
	if !contains(keys, c.Version) {
		return fmt.Errorf("appgate client version invalid, got %d, default is %d", c.Version, DefaultClientVersion)
	}
	return nil
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

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok, err := FileExists(c.PemFilePath); err == nil && ok {
		certs, err := os.ReadFile(c.PemFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not read pem file %w", err)
		}
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			return nil, fmt.Errorf("unable to append cert %s", c.PemFilePath)
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.Insecure,
			RootCAs:            rootCAs,
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

func guessVersion(clientVersion int) (*version.Version, error) {
	// TODO query GET /appliance controller and check exact version.
	// POST /login does not include version anymore.
	switch clientVersion {
	case Version13:
		return version.NewVersion("5.2.0+estimated")
	case Version14:
		return version.NewVersion("5.3.0+estimated")
	case Version15:
		return version.NewVersion("5.4.0+estimated")
	case Version16:
		return version.NewVersion("5.5.0+estimated")

	}
	return nil, fmt.Errorf("could not determine appliance version with client version %d", clientVersion)
}

// GetToken makes first login and initiate the client towards the controller.
// this is always the first request made
func (c *Client) GetToken() (string, error) {
	if len(c.Config.BearerToken) > 0 {
		log.Printf("[DEBUG] Authenticate with Bearer token provided as APPGATE_BEARER_TOKEN")
		c.Token = fmt.Sprintf("Bearer %s", c.Config.BearerToken)
	}
	cfg := c.Config
	latestSupportedVersion, err := version.NewVersion(ApplianceVersionMap[DefaultClientVersion])
	if err != nil {
		return "", err
	}
	currentVersion, err := guessVersion(cfg.Version)
	if err != nil {
		return "", err
	}
	c.ApplianceVersion = currentVersion
	c.LatestSupportedVersion = latestSupportedVersion

	if len(c.Token) > 0 {
		return c.Token, nil
	}
	response, err := c.login()
	if err != nil {
		return "", err
	}

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
	// if the appgatesdp provider is combined with, for example  aws_instance where we create the initial controller
	// it might take awhile for the controller to startup and be responsive, so until its up it can return 500, 502, 503
	// these status code is treated as retryable errors, during exponentialBackOff.MaxElapsedTime window.
	// we will use this exponential backoff to retry until we get a 200-400 HTTP response from /login
	exponentialBackOff.MaxElapsedTime = c.loginTimoutDuration()
	loginResponse := &openapi.LoginResponse{}
	err := backoff.Retry(func() error {
		login, response, err := c.API.LoginApi.LoginPost(ctx).LoginRequest(loginOpts).Execute()
		if response == nil {
			if err != nil {
				if err, ok := err.(*url.Error); ok {
					if err, ok := err.Unwrap().(x509.UnknownAuthorityError); ok {
						return &backoff.PermanentError{
							Err: fmt.Errorf("Import certificate or toggle APPGATE_INSECURE - %s", err),
						}
					}
				}
			}
			log.Printf("[DEBUG] Login failed, No response %s", err)
			return fmt.Errorf("No response from controller %w", err)
		}
		if response.StatusCode >= 500 {
			log.Printf("[DEBUG] Login failed, controller not responding, got HTTP %d", response.StatusCode)
			return fmt.Errorf("Controller got %w", err)
		}
		if err != nil {
			if err, ok := err.(openapi.GenericOpenAPIError); ok {
				if err, ok := err.Model().(openapi.InlineResponse406); ok {
					return &backoff.PermanentError{
						Err: fmt.Errorf(
							"You are using the wrong client_version (peer api version) for you appgate sdp collective, you are using %d; min: %d max: %d",
							c.Config.Version,
							err.GetMinSupportedVersion(),
							err.GetMaxSupportedVersion(),
						),
					}
				}
			}
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
