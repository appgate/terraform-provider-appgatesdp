package provider

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

const testDeviceID = "4c07bc67-57ea-42dd-b702-c2d6c45419fc"

var protoV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
	"appgatesdp": providerserver.NewProtocol6WithError(New("test")()),
}

// terraformSlots bounds how many per-resource subtests drive Terraform at
// once. Each step starts Terraform and a provider server; too many at once
// makes provider startup time out ("timeout waiting on reattach config").
var terraformSlots = make(chan struct{}, max(2, runtime.NumCPU()/3))

// parallelTerraform runs a subtest in parallel within the Terraform limit.
func parallelTerraform(t *testing.T) {
	t.Helper()
	t.Parallel()
	terraformSlots <- struct{}{}
	t.Cleanup(func() { <-terraformSlots })
}

// forEachVersion runs fn once per registered API version.
func forEachVersion(t *testing.T, fn func(t *testing.T, v apiversion.Version)) {
	t.Helper()
	for _, v := range apiversion.Default().Versions {
		t.Run(v.APIVersion, func(t *testing.T) {
			t.Parallel()
			fn(t, v)
		})
	}
}

// newServer starts a mock controller speaking API version v.
func newServer(t *testing.T, v string) *sdptest.Server {
	t.Helper()
	api, err := generated.API(v)
	if err != nil {
		t.Fatal(err)
	}
	return sdptest.New(t, api)
}

// providerConfig returns a provider block pointing at the mock controller.
func providerConfig(srv *sdptest.Server, version string) string {
	apiVersion := ""
	if version != "" {
		apiVersion = fmt.Sprintf("api_version = %q", version)
	}
	return fmt.Sprintf(`
provider "appgatesdp" {
  url       = %q
  username  = %q
  password  = %q
  device_id = %q
  %s
}
`, srv.URL, sdptest.Username, sdptest.Password, testDeviceID, apiVersion)
}

// checkAcceptHeaders verifies every request used the version's media type.
func checkAcceptHeaders(srv *sdptest.Server, v apiversion.Version) func(*terraform.State) error {
	return func(*terraform.State) error {
		for _, r := range srv.Requests() {
			if r.Accept != v.AcceptHeader() && !srv.AcceptOverride(r.Path, r.Accept) {
				return fmt.Errorf("%s %s sent Accept %q, want %q", r.Method, r.Path, r.Accept, v.AcceptHeader())
			}
		}
		return nil
	}
}

// lastRequest returns the body of the last request with method and path prefix.
func lastRequest(srv *sdptest.Server, method, pathPrefix string) map[string]any {
	reqs := srv.Requests()
	slices.Reverse(reqs)
	for _, r := range reqs {
		if r.Method == method && strings.HasPrefix(r.Path, pathPrefix) {
			return r.Body
		}
	}
	return nil
}

func resourceID(s *terraform.State, name string) (string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return "", fmt.Errorf("%s not in state", name)
	}
	return rs.Primary.ID, nil
}

func newTestClient(t *testing.T, srv *sdptest.Server, version int) *sdp.Client {
	t.Helper()
	c, err := sdp.New(sdp.Config{URL: srv.URL, Username: sdptest.Username, Password: sdptest.Password, DeviceID: testDeviceID}, version)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
