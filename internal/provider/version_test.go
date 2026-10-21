package provider

import (
	"fmt"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

// supportedIn reports whether version v is one of versions.
func supportedIn(v apiversion.Version, versions ...string) bool {
	return slices.Contains(versions, v.APIVersion)
}

// TestResourceAvailabilityByVersion: resources missing from a version fail at
// plan time with an actionable message instead of calling the API.
func TestResourceAvailabilityByVersion(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		step := resource.TestStep{
			Config: providerConfig(srv, v.APIVersion) + `
resource "appgatesdp_replication_target" "test" {
  name             = "tf-target"
  replication_tags = ["replicate"]
}
`,
		}
		if supportedIn(v, "v24", "v25") {
			step.Check = resource.TestCheckResourceAttrSet("appgatesdp_replication_target.test", "id")
		} else {
			step.ExpectError = regexp.MustCompile(`appgatesdp_replication_target is not available in SDP Admin API ` + v.APIVersion)
		}
		resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: protoV6Factories, Steps: []resource.TestStep{step}})
		if !supportedIn(v, "v24", "v25") {
			for _, r := range srv.Requests() {
				if r.Path != "/login" {
					t.Errorf("unsupported resource must not reach the API, got %s %s", r.Method, r.Path)
				}
			}
		}
	})
}

// TestAttributeAvailabilityByVersion covers attributes added (policy
// block_local_dns_requests, v23) and removed (entitlement risk_sensitivity,
// gone in v25) between versions.
func TestAttributeAvailabilityByVersion(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		cfg := providerConfig(srv, v.APIVersion) + `
resource "appgatesdp_access_policy" "test" {
  name                     = "tf-policy"
  expression               = "return true;"
  block_local_dns_requests = true
}
`
		step := resource.TestStep{Config: cfg}
		if supportedIn(v, "v23", "v24", "v25") {
			step.Check = func(*terraform.State) error {
				body := lastRequest(srv, "POST", "/policies")
				if body["blockLocalDnsRequests"] != true || body["type"] != "Access" {
					return fmt.Errorf("unexpected POST body %v", body)
				}
				return nil
			}
		} else {
			step.ExpectError = regexp.MustCompile(`(?s)block_local_dns_requests is not available in SDP\s+Admin API ` + v.APIVersion)
		}
		resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: protoV6Factories, Steps: []resource.TestStep{step}})
	})
}

func TestRemovedAttribute(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		site := srv.Seed("/sites", map[string]any{"name": "site"})
		step := resource.TestStep{Config: providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_entitlement" "test" {
  name             = "tf-entitlement"
  site             = %q
  risk_sensitivity = "High"
  conditions       = []
  actions = [{
    subtype = "tcp_up"
    action  = "allow"
    hosts   = ["10.0.0.1"]
    ports   = ["443"]
  }]
}
`, site)}
		if v.APIVersion == "v25" {
			step.ExpectError = regexp.MustCompile(`(?s)risk_sensitivity is not available in SDP\s+Admin API v25`)
		} else {
			step.Check = func(*terraform.State) error {
				if body := lastRequest(srv, "POST", "/entitlements"); body["riskSensitivity"] != "High" {
					return fmt.Errorf("riskSensitivity not sent: %v", body)
				}
				return nil
			}
		}
		resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: protoV6Factories, Steps: []resource.TestStep{step}})
	})
}

// TestVersionNegotiation: without api_version the newest version supported
// by both sides is used for every request.
func TestVersionNegotiation(t *testing.T) {
	srv := newServer(t, "v23")
	srv.MinVersion, srv.MaxVersion = 18, 23
	want, _ := apiversion.Default().Lookup("v23")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "") + `
resource "appgatesdp_criteria_script" "test" {
  name       = "tf-script"
  expression = "return true;"
}
`,
			Check: func(*terraform.State) error {
				for _, r := range srv.Requests() {
					if r.Path == "/login" && r.Body["username"] == nil {
						continue // negotiation probe
					}
					if r.Accept != want.AcceptHeader() {
						return fmt.Errorf("%s %s used %q", r.Method, r.Path, r.Accept)
					}
				}
				return nil
			},
		}},
	})
}

func TestVersionNegotiationNoOverlap(t *testing.T) {
	srv := newServer(t, "v21")
	srv.MinVersion, srv.MaxVersion = 10, 15
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config:      providerConfig(srv, "") + `data "appgatesdp_global_settings" "test" {}`,
			ExpectError: regexp.MustCompile(`controller supports API v10 to v15`),
		}},
	})
}

func TestUnsupportedAPIVersion(t *testing.T) {
	srv := newServer(t, "v25")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config:      providerConfig(srv, "v99") + `data "appgatesdp_global_settings" "test" {}`,
			ExpectError: regexp.MustCompile(`unsupported API version "v99"; supported versions: v21, v22, v23, v24, v25`),
		}},
	})
}

func TestLegacyClientVersion(t *testing.T) {
	srv := newServer(t, "v22")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
provider "appgatesdp" {
  url            = %q
  username       = %q
  password       = %q
  client_version = 22
}
data "appgatesdp_global_settings" "test" {}
`, srv.URL, sdptest.Username, sdptest.Password),
			Check: func(*terraform.State) error {
				for _, r := range srv.Requests() {
					if r.Accept != "application/vnd.appgate.peer-v22+json" {
						return fmt.Errorf("%s used %q", r.Path, r.Accept)
					}
				}
				return nil
			},
		}},
	})
}
