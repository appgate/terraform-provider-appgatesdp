// Package acctest connects acceptance tests to real SDP environments.
//
// Environments are configured per API version through environment variables:
//
//	SDP_ACC_V25_URL        controller admin URL (required to enable v25)
//	SDP_ACC_V25_USERNAME   admin username
//	SDP_ACC_V25_PASSWORD   admin password
//	SDP_ACC_V25_PROVIDER   identity provider (default "local")
//	SDP_ACC_V25_PEM        CA certificate file; without it TLS is not verified
//
// Tests only run with TF_ACC=1. Tests that create or delete objects also
// require SDP_ACC_ALLOW_DESTRUCTIVE=1 and only touch objects named with
// NamePrefix, which the sweeper removes.
package acctest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

// NamePrefix marks objects created by acceptance tests. IdentifierPrefix is
// used where names must be JavaScript identifiers (e.g. criteria scripts).
const (
	NamePrefix       = "tf-acc-"
	IdentifierPrefix = "tf_acc_"
)

// Env is one configured SDP environment.
type Env struct {
	Version  apiversion.Version
	URL      string
	Username string
	Password string
	Provider string
	PEMFile  string
}

// Environments returns the environments configured for registered versions.
func Environments() []Env {
	var out []Env
	for _, v := range apiversion.Default().Versions {
		key := "SDP_ACC_" + strings.ToUpper(v.APIVersion) + "_"
		if os.Getenv(key+"URL") == "" {
			continue
		}
		out = append(out, Env{
			Version:  v,
			URL:      os.Getenv(key + "URL"),
			Username: os.Getenv(key + "USERNAME"),
			Password: os.Getenv(key + "PASSWORD"),
			Provider: envOr(key+"PROVIDER", "local"),
			PEMFile:  os.Getenv(key + "PEM"),
		})
	}
	return out
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ForEachEnv runs fn for every configured environment. It skips unless
// TF_ACC is set and at least one environment is configured.
func ForEachEnv(t *testing.T, fn func(t *testing.T, env Env)) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance tests require TF_ACC=1")
	}
	envs := Environments()
	if len(envs) == 0 {
		t.Skip("no SDP environment configured; set SDP_ACC_<VERSION>_URL, e.g. SDP_ACC_V25_URL")
	}
	for _, env := range envs {
		t.Run(env.Version.APIVersion, func(t *testing.T) { fn(t, env) })
	}
}

// RequireDestructive skips tests that change the environment unless they are
// explicitly allowed.
func RequireDestructive(t *testing.T) {
	t.Helper()
	if os.Getenv("SDP_ACC_ALLOW_DESTRUCTIVE") != "1" {
		t.Skip("test creates and deletes objects; set SDP_ACC_ALLOW_DESTRUCTIVE=1 to run it against the configured environments")
	}
}

// Name returns a unique object name with NamePrefix.
func Name(kind string) string {
	return NamePrefix + kind + "-" + uuid.NewString()[:8]
}

// Identifier returns a unique JavaScript identifier with IdentifierPrefix.
func Identifier(kind string) string {
	return IdentifierPrefix + strings.ReplaceAll(kind, "-", "_") + "_" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
}

// isTestObject reports whether a name or identifier was created by tests.
func isTestObject(s string) bool {
	return strings.HasPrefix(s, NamePrefix) || strings.HasPrefix(s, IdentifierPrefix) || strings.Contains(s, NamePrefix)
}

// SkipUnlicensed is a TestCase.ErrorCheck that skips the test when the
// collective's license does not include the feature (HTTP 402); that is a
// property of the environment, not of the provider.
func SkipUnlicensed(t *testing.T) func(error) error {
	return func(err error) error {
		if err != nil && strings.Contains(err.Error(), "HTTP 402 insufficient license") {
			t.Skipf("feature not licensed in this environment: %v", err)
		}
		return err
	}
}

// ProviderConfig renders the provider block of an environment.
func (e Env) ProviderConfig() string {
	pem := ""
	if e.PEMFile != "" {
		pem = fmt.Sprintf("pem_filepath = %q\n  insecure = false", e.PEMFile)
	}
	return fmt.Sprintf(`
provider "appgatesdp" {
  url         = %q
  username    = %q
  password    = %q
  provider    = %q
  api_version = %q
  %s
}
`, e.URL, e.Username, e.Password, e.Provider, e.Version.APIVersion, pem)
}

// Client returns an API client for the environment.
func (e Env) Client(t *testing.T) (*sdp.Client, *apimeta.API) {
	t.Helper()
	c, err := sdp.New(sdp.Config{
		URL: e.URL, Username: e.Username, Password: e.Password, Provider: e.Provider,
		PEMFile: e.PEMFile, Insecure: e.PEMFile == "", DeviceID: uuid.NewString(), LoginTimeout: time.Minute,
	}, e.Version.Number())
	if err != nil {
		t.Fatal(err)
	}
	api, err := generated.API(e.Version.APIVersion)
	if err != nil {
		t.Fatal(err)
	}
	return c, api
}

// Sweep deletes every object named with NamePrefix from the collections of
// the environment's API version.
func Sweep(ctx context.Context, c *sdp.Client, api *apimeta.API) error {
	seen := map[string]bool{}
	var errs []string
	for _, name := range sortedKeys(api.Resources) {
		r := api.Resources[name]
		if r.Singleton || r.Adopt || seen[r.Path] {
			continue
		}
		seen[r.Path] = true
		q := url.Values{"query": {NamePrefix}, "range": {"0-1000"}}
		if r.ReadList != "" {
			// Not every list endpoint searches the identifier; filter locally.
			q = nil
		}
		var list map[string]any
		if r.ReadList == "" {
			// The search query matches names; test identifiers use a second prefix.
			q.Set("query", "tf")
		}
		if err := c.Do(ctx, http.MethodGet, r.Path, q, nil, &list); err != nil {
			// Unlicensed features (HTTP 402) cannot hold test objects.
			var apiErr *sdp.Error
			if !errors.As(err, &apiErr) || apiErr.Status != http.StatusPaymentRequired {
				errs = append(errs, err.Error())
			}
			continue
		}
		items, _ := list[r.ListProperty()].([]any)
		for _, item := range items {
			obj, _ := item.(map[string]any)
			name, _ := obj["name"].(string)
			id := fmt.Sprint(obj[r.ID()])
			// Objects without a name (e.g. blacklist entries) carry the
			// prefix in their identifier.
			if !isTestObject(name) && !isTestObject(id) {
				continue
			}
			if err := c.Do(ctx, http.MethodDelete, r.DeleteObjectPath(url.PathEscape(id)), nil, nil, nil); err != nil && !sdp.IsNotFound(err) {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("sweep: %s", strings.Join(errs, "; "))
	}
	return nil
}
