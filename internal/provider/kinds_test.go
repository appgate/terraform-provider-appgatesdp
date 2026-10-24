package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

// TestSingletonResource: global settings are updated with PUT and reset with
// DELETE; they are never created.
func TestSingletonResource(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if len(srv.Object("/global-settings", "")) != 0 {
					return fmt.Errorf("destroy did not reset global settings")
				}
				return nil
			},
			Steps: []resource.TestStep{{
				Config: providerConfig(srv, v.APIVersion) + `
resource "appgatesdp_global_settings" "test" {
  administration_token_expiration = 60
  audit_log_persistence_mode      = "Default"
  claims_token_expiration         = 1440
  entitlement_token_expiration    = 180
  vpn_certificate_expiration      = 525600
  backup_passphrase               = "s3cret"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("appgatesdp_global_settings.test", "id", "global_settings"),
					resource.TestCheckResourceAttr("appgatesdp_global_settings.test", "claims_token_expiration", "1440"),
					resource.TestCheckResourceAttr("appgatesdp_global_settings.test", "backup_passphrase", "s3cret"),
					func(*terraform.State) error {
						for _, r := range srv.Requests() {
							if r.Method == http.MethodPost && r.Path != "/login" {
								return fmt.Errorf("singleton must not POST, got %s", r.Path)
							}
						}
						if got := srv.Object("/global-settings", "")["claimsTokenExpiration"]; fmt.Sprint(got) != "1440" {
							return fmt.Errorf("claimsTokenExpiration = %v", got)
						}
						return nil
					},
				),
			}},
		})
	})
}

// TestWriteOnlySecret: the controller never returns the LDAP admin password;
// it is kept in state, resent on update and causes no diff on refresh.
func TestWriteOnlySecret(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		cfg := func(port int) string {
			return providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_ldap_identity_provider" "test" {
  name                     = "tf-ldap"
  hostnames                = ["ldap.example.com"]
  port                     = %d
  admin_distinguished_name = "CN=admin,DC=example,DC=com"
  admin_password           = "hunter2"
}
`, port)
		}
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			Steps: []resource.TestStep{
				{
					Config: cfg(389),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("appgatesdp_ldap_identity_provider.test", "admin_password", "hunter2"),
						func(*terraform.State) error {
							if body := lastRequest(srv, "POST", "/identity-providers"); body["type"] != "Ldap" || body["adminPassword"] != "hunter2" {
								return fmt.Errorf("unexpected POST body %v", body)
							}
							return nil
						},
					),
				},
				{
					Config: cfg(636),
					Check: func(*terraform.State) error {
						body := lastRequest(srv, "PUT", "/identity-providers/")
						if body["adminPassword"] != "hunter2" || fmt.Sprint(body["port"]) != "636" {
							return fmt.Errorf("update must resend the write-only password, got %v", body)
						}
						return nil
					},
				},
			},
		})
	})
}

// TestVariantTypeChange: an object whose discriminator changed outside
// Terraform is no longer the resource and is planned for re-creation.
func TestVariantTypeChange(t *testing.T) {
	srv := newServer(t, "v25")
	cfg := providerConfig(srv, "v25") + `
resource "appgatesdp_dns_policy" "test" {
  name       = "tf-dns"
  expression = "return true;"
}
`
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: func(s *terraform.State) (err error) {
					id, err = resourceID(s, "appgatesdp_dns_policy.test")
					return err
				},
			},
			{
				PreConfig:          func() { srv.Update("/policies", id, func(o map[string]any) { o["type"] = "Access" }) },
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAdoptBuiltInObject: the local identity provider always exists; create
// adopts it by name and destroy leaves it in place.
func TestAdoptBuiltInObject(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		local := srv.Seed("/identity-providers", map[string]any{"name": "local", "type": "LocalDatabase"})
		srv.Seed("/identity-providers", map[string]any{"name": "local", "type": "Ldap"})
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if srv.Object("/identity-providers", local) == nil {
					return fmt.Errorf("built-in identity provider was deleted")
				}
				return nil
			},
			Steps: []resource.TestStep{{
				Config: providerConfig(srv, v.APIVersion) + `
resource "appgatesdp_local_database_identity_provider" "test" {
  name  = "local"
  notes = "managed"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("appgatesdp_local_database_identity_provider.test", "id", local),
					func(*terraform.State) error {
						if got := srv.Object("/identity-providers", local)["notes"]; got != "managed" {
							return fmt.Errorf("notes = %v", got)
						}
						return nil
					},
				),
			}},
		})
	})
}

// TestUserChosenID: SDP accepts client generated UUIDs; changing the id
// replaces the object.
func TestUserChosenID(t *testing.T) {
	srv := newServer(t, "v25")
	cfg := func(id string) string {
		return providerConfig(srv, "v25") + fmt.Sprintf(`
resource "appgatesdp_ip_pool" "test" {
  id   = %q
  name = "tf-pool"
}
`, id)
	}
	const id1, id2 = "a1b2c3d4-0000-4000-8000-000000000001", "a1b2c3d4-0000-4000-8000-000000000002"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{
			{Config: cfg(id1), Check: resource.TestCheckResourceAttr("appgatesdp_ip_pool.test", "id", id1)},
			{Config: cfg(id2), Check: func(*terraform.State) error {
				if srv.Object("/ip-pools", id1) != nil || srv.Object("/ip-pools", id2) == nil {
					return fmt.Errorf("changing id must replace the object")
				}
				return nil
			}},
		},
	})
}

func TestAPIErrors(t *testing.T) {
	tests := map[string]struct {
		setup  func(srv *sdptest.Server)
		config func(srv *sdptest.Server) string
		want   string
	}{
		"validation error": {
			setup: func(srv *sdptest.Server) {
				srv.Fail("POST", "/conditions", 422, `{"id":"validation error","message":"Validation error.","errors":[{"field":"expression","message":"syntax error"}]}`, 1)
			},
			want: `(?s)HTTP 422 validation error: Validation error.\s+expression: syntax error`,
		},
		"server error": {
			setup: func(srv *sdptest.Server) {
				srv.Fail("POST", "/conditions", 500, `{"id":"unexpected","message":"boom"}`, 1)
			},
			want: `HTTP 500 unexpected: boom`,
		},
		"authentication failure": {
			config: func(srv *sdptest.Server) string {
				return fmt.Sprintf(`
provider "appgatesdp" {
  url         = %q
  username    = "admin"
  password    = "wrong"
  api_version = "v25"
}
`, srv.URL)
			},
			want: `login failed: POST /login: HTTP 401 unauthorized`,
		},
		"unexpected response": {
			setup: func(srv *sdptest.Server) {
				srv.Seed("/conditions", map[string]any{"id": "x", "name": []any{"not", "a", "string"}})
				srv.Fail("POST", "/conditions", 200, `{"id":"x"}`, 1)
			},
			want: `Unexpected API response for appgatesdp_condition`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv := newServer(t, "v25")
			if tc.setup != nil {
				tc.setup(srv)
			}
			cfg := providerConfig(srv, "v25")
			if tc.config != nil {
				cfg = tc.config(srv)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: protoV6Factories,
				Steps: []resource.TestStep{{
					Config: cfg + `
resource "appgatesdp_condition" "test" {
  name       = "tf-condition"
  expression = "return true;"
}
`,
					ExpectError: regexp.MustCompile(tc.want),
				}},
			})
		})
	}
}

func TestDataSources(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		id := srv.Seed("/conditions", map[string]any{"name": "always", "expression": "return true;", "tags": []any{"a"}})
		srv.Seed("/conditions", map[string]any{"name": "always-not", "expression": "return false;"})
		srv.Seed("/conditions", map[string]any{"name": "twin", "expression": "1"})
		srv.Seed("/conditions", map[string]any{"name": "twin", "expression": "2"})
		base := providerConfig(srv, v.APIVersion)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			Steps: []resource.TestStep{
				{
					Config: base + fmt.Sprintf(`
data "appgatesdp_condition" "by_id" { id = %q }
data "appgatesdp_condition" "by_name" { name = "always" }
`, id),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.appgatesdp_condition.by_id", "name", "always"),
						resource.TestCheckResourceAttr("data.appgatesdp_condition.by_id", "tags.0", "a"),
						resource.TestCheckResourceAttr("data.appgatesdp_condition.by_name", "id", id),
						resource.TestCheckResourceAttr("data.appgatesdp_condition.by_name", "expression", "return true;"),
					),
				},
				{Config: base + `data "appgatesdp_condition" "x" { name = "twin" }`, ExpectError: regexp.MustCompile(`2 objects named "twin"`)},
				{Config: base + `data "appgatesdp_condition" "x" { name = "missing" }`, ExpectError: regexp.MustCompile(`no object named "missing"`)},
				{Config: base + `data "appgatesdp_condition" "x" {}`, ExpectError: regexp.MustCompile(`Set either id or name`)},
				{Config: base + `data "appgatesdp_condition" "x" { id = "00000000-0000-4000-8000-000000000000" }`, ExpectError: regexp.MustCompile(`HTTP 404`)},
			},
		})
	})
}

// TestSingletonCreatedOnFirstWrite: a singleton that does not exist until
// first written (GET 404, e.g. /my-preferences) is created from the plan.
func TestSingletonCreatedOnFirstWrite(t *testing.T) {
	srv := newServer(t, "v25")
	srv.Fail(http.MethodGet, "/my-preferences", http.StatusNotFound, `{"id":"not found","message":"Not found."}`, 1)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "v25") + `
resource "appgatesdp_my_preferences" "test" {
  value = jsonencode({ theme = "dark" })
}
`,
			Check: func(*terraform.State) error {
				if body := lastRequest(srv, http.MethodPut, "/my-preferences"); body["theme"] != "dark" {
					return fmt.Errorf("PUT body %v", body)
				}
				return nil
			},
		}},
	})
}
