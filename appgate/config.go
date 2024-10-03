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
	"sync"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v21/openapi"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-version"
	"golang.org/x/net/http/httpproxy"
)

const (
	// DefaultDescription is the default string for terraform resources.
	DefaultDescription = "Managed by terraform"
)

// Config for appgate provider.
type Config struct {
	URL          string        `json:"appgate_url,omitempty"`
	Username     string        `json:"appgate_username,omitempty"`
	Password     string        `json:"appgate_password,omitempty"`
	Provider     string        `json:"appgate_provider,omitempty"`
	Insecure     bool          `json:"appgate_insecure,omitempty"`
	Timeout      int           `json:"appgate_timeout,omitempty"`
	LoginTimeout time.Duration `json:"appgate_login_timeout,omitempty"`
	Debug        bool          `json:"appgate_http_debug,omitempty"`
	Version      int           `json:"appgate_client_version,omitempty"`
	BearerToken  string        `json:"appgate_bearer_token,omitempty"`
	PemFilePath  string        `json:"appgate_pem_filepath,omitempty"`
	DeviceID     string        `json:"appgate_device_id,omitempty"`
	UserAgent    string
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
			return fmt.Errorf("appgate bearer_token set, but invalid format, expected base64 %w", err)
		}
	} else if len(c.Username) < 1 && len(c.Password) < 1 {
		return fmt.Errorf("username and password required if appgate bearer token is empty")
	}

	return nil
}

// Client is the appgate API client.
type Client struct {
	mu               sync.Mutex
	Token            string
	UUID             string
	ApplianceVersion *version.Version
	ClientVersion    int
	API              *openapi.APIClient
	Config           *Config
}

var proxyFunc func(*url.URL) (*url.URL, error)

func proxyFromEnvironment(req *http.Request) (*url.URL, error) {
	if key, ok := os.LookupEnv("HTTP_PROXY"); ok {
		proxyURL, err := url.Parse(key)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Using HTTP PROXY %s", proxyURL)
		proxyConfig := &httpproxy.Config{
			// use the same Addr for both HTTPS_PROXY and HTTP_PROXY for backwards compatibility
			HTTPProxy:  proxyURL.String(),
			HTTPSProxy: proxyURL.String(),
		}
		if noProxy, ok := os.LookupEnv("NO_PROXY"); ok {
			log.Printf("[DEBUG] Using NO_PROXY %s", noProxy)
			proxyConfig.NoProxy = noProxy
		}
		proxyFunc = proxyConfig.ProxyFunc()
	}
	if proxyFunc == nil {
		proxyFunc = httpproxy.FromEnvironment().ProxyFunc()
		log.Printf("[DEBUG] Using HTTP PROXY FromEnvironment")
	}
	return proxyFunc(req.URL)
}

// Client creates the http client, APIClient, and setup configuration for
// custom pem file
// toggle tls verification based on config
// setup http proxy based on environment variables
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
		Proxy:               proxyFromEnvironment,
	}

	httpclient := &http.Client{
		Transport: tr,
		Timeout:   ((timeoutDuration * 2) * time.Second),
	}

	clientCfg := &openapi.Configuration{
		DefaultHeader: map[string]string{
			"Accept": fmt.Sprintf("application/vnd.appgate.peer-v%d+json", c.Version),
		},
		UserAgent: c.UserAgent,
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
	case Version18:
		return version.NewVersion("6.1.0+estimated")
	case Version19:
		return version.NewVersion("6.2.0+estimated")
	case Version20:
		return version.NewVersion("6.3.0+estimated")
	case Version21:
		return version.NewVersion("6.4.0+estimated")
	}
	return nil, fmt.Errorf("could not determine appliance version with client version %d", clientVersion)
}

// GetToken makes first login and initiate the client towards the controller.
// this is always the first request made
func (c *Client) GetToken() (string, error) {
	// we will do a lock here to avoid concurrent race condition if we need to update the root
	// client default header values and only authenticate if we haven't already cached a bearer token.
	c.mu.Lock()
	defer c.mu.Unlock()
	cfg := c.Config

	if len(cfg.BearerToken) > 0 {
		log.Printf("[DEBUG] Authenticate with Bearer token provided as APPGATE_BEARER_TOKEN")
		c.Token = fmt.Sprintf("Bearer %s", cfg.BearerToken)
		return c.Token, nil
	}
	if len(c.Token) > 0 {
		log.Printf("[DEBUG] Using existing token")
		return c.Token, nil
	}

	// if the client_version is set to the default minimum value, we will do
	// a error request to login to determine the maximum allowed version for the current
	// controller to use.
	// This only happens if the provisioner has omitted the client_version from their provider configuration.
	if cfg.Version == MinimumSupportedVersion {
		var minMaxErr *minMaxError
		_, err := c.login(context.WithValue(
			context.Background(),
			openapi.ContextAcceptHeader,
			fmt.Sprintf("application/vnd.appgate.peer-v%d+json", 5),
		))
		if errors.As(err, &minMaxErr) {
			log.Printf("[DEBUG] retrieved client version %d to use from login error response", minMaxErr.Max)
			cfg.Version = int(minMaxErr.Max)
		} else {
			log.Printf("[DEBUG] could not compute client version API support, fallback %d", DefaultClientVersion)
			cfg.Version = DefaultClientVersion
		}
		c.API.GetConfig().DefaultHeader["Accept"] = fmt.Sprintf("application/vnd.appgate.peer-v%d+json", cfg.Version)
	}

	currentVersion, err := guessVersion(cfg.Version)
	if err != nil {
		return "", err
	}
	c.ApplianceVersion = currentVersion

	if len(c.Token) > 0 {
		return c.Token, nil
	}
	response, err := c.login(context.Background())
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

type minMaxError struct {
	Err      error
	Min, Max int32
}

func (e *minMaxError) Error() string {
	return e.Err.Error()
}

func (c *Client) login(ctx context.Context) (*openapi.LoginResponse, error) {
	loginOpts := openapi.LoginRequest{
		ProviderName: c.Config.Provider,
		Username:     openapi.PtrString(c.Config.Username),
		Password:     openapi.PtrString(c.Config.Password),
		DeviceId:     c.Config.DeviceID,
	}

	// Since /login is the first request we do, it provide us the earliest check if a controller is up and running
	// if the appgatesdp provider is combined with, for example  aws_instance where we create the initial controller
	// it might take awhile for the controller to startup and be responsive, so until its up it can return 500, 502, 503
	// these status code is treated as retryable errors, during exponentialBackOff.MaxElapsedTime window.
	// we will use this exponential backoff to retry until we get a 200-400 HTTP response from /login
	exponentialBackOff.MaxElapsedTime = c.Config.LoginTimeout
	loginResponse := &openapi.LoginResponse{}
	err := backoff.Retry(func() error {
		login, response, err := c.API.LoginApi.LoginPost(ctx).LoginRequest(loginOpts).Execute()
		if response == nil {
			if err != nil && errors.As(err, &x509.UnknownAuthorityError{}) {
				return &backoff.PermanentError{
					Err: fmt.Errorf("Import certificate or toggle APPGATE_INSECURE - %s", err),
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
			if err, ok := err.(*openapi.GenericOpenAPIError); ok {
				if responseErr, ok := err.Model().(openapi.LoginPost406Response); ok {
					return &backoff.PermanentError{Err: &minMaxError{
						Err: fmt.Errorf("Invalid appgatesdp.client_version for your collective %w", err),
						Min: responseErr.GetMinSupportedVersion(),
						Max: responseErr.GetMaxSupportedVersion(),
					}}
				}
			}
			log.Printf("[DEBUG] Login failed permanently, got HTTP %d", response.StatusCode)
			return &backoff.PermanentError{Err: err}
		}
		loginResponse = login
		return nil
	}, backoff.WithContext(&exponentialBackOff, ctx))
	if err != nil {
		return nil, prettyPrintAPIError(err)
	}
	log.Printf("[DEBUG] Login OK")
	return loginResponse, nil
}
