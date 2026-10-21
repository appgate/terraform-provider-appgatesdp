package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
)

func init() {
	// Keep appliance state polling fast in tests.
	statusPollInterval = 10 * time.Millisecond
}

// seededAppliance is a complete appliance as returned by the controller.
func seededAppliance(name string, activated bool) map[string]any {
	return map[string]any{
		"name": name, "hostname": name + ".example.com", "activated": activated,
		"clientInterface": map[string]any{"hostname": name + ".example.com"},
		"networking":      map[string]any{},
		"controller":      map[string]any{"enabled": false},
	}
}

func expectReplace(addr string) resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionReplace),
	}}
}

// TestLicense: licenses are uploaded with POST /license, read from the
// licenses list of GET /license and deleted by id; any change replaces them.
func TestLicense(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		cfg := func(lic string) string {
			return providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_license" "test" {
  license = %q
}
`, lic)
		}
		var firstID string
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if n := srv.Count("/license"); n != 0 {
					return fmt.Errorf("%d licenses left", n)
				}
				return nil
			},
			Steps: []resource.TestStep{
				{
					Config: cfg("license-1"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("appgatesdp_license.test", "id"),
						resource.TestCheckResourceAttr("appgatesdp_license.test", "license", "license-1"),
						func(s *terraform.State) (err error) {
							firstID, err = resourceID(s, "appgatesdp_license.test")
							if body := lastRequest(srv, http.MethodPost, "/license"); body["license"] != "license-1" {
								return fmt.Errorf("POST body %v", body)
							}
							return err
						},
					),
				},
				{
					Config:           cfg("license-2"),
					ConfigPlanChecks: expectReplace("appgatesdp_license.test"),
					Check: func(s *terraform.State) error {
						id, _ := resourceID(s, "appgatesdp_license.test")
						if id == firstID || srv.Object("/license", firstID) != nil || srv.Count("/license") != 1 {
							return fmt.Errorf("changing the license must replace it")
						}
						return nil
					},
				},
				{ResourceName: "appgatesdp_license.test", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"license"}},
				{
					// Removed outside Terraform: the license is uploaded again.
					PreConfig: func() {
						for _, id := range srv.IDs("/license") {
							srv.Delete("/license", id)
						}
					},
					Config: cfg("license-2"),
					Check: func(*terraform.State) error {
						if srv.Count("/license") != 1 {
							return fmt.Errorf("license was not uploaded again")
						}
						return nil
					},
				},
			},
		})
	})
}

// TestBlacklistUser: the distinguished name is the id, any change replaces
// the entry and destroy removes it with DELETE /blacklist/{dn}.
func TestBlacklistUser(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		const dn = "CN=alice,OU=ldap"
		cfg := func(reason string) string {
			return providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_blacklist_user" "test" {
  user_distinguished_name = %q
  reason                  = %q
}
`, dn, reason)
		}
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if srv.Object("/blacklist", dn) != nil {
					return fmt.Errorf("user is still blacklisted")
				}
				for _, r := range srv.Requests() {
					if r.Method == http.MethodDelete && r.Path == "/blacklist/"+dn {
						return nil
					}
				}
				return fmt.Errorf("DELETE /blacklist/%s was not called", dn)
			},
			Steps: []resource.TestStep{
				{
					Config: cfg("left the company"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("appgatesdp_blacklist_user.test", "id", dn),
						resource.TestCheckResourceAttr("appgatesdp_blacklist_user.test", "reason", "left the company"),
					),
				},
				{
					Config:           cfg("compromised"),
					ConfigPlanChecks: expectReplace("appgatesdp_blacklist_user.test"),
					Check:            resource.TestCheckResourceAttr("appgatesdp_blacklist_user.test", "reason", "compromised"),
				},
				{ResourceName: "appgatesdp_blacklist_user.test", ImportState: true, ImportStateId: dn, ImportStateVerify: true},
			},
		})
	})
}

func TestBlacklistUserRequiresDistinguishedName(t *testing.T) {
	srv := newServer(t, "v25")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config:      providerConfig(srv, "v25") + `resource "appgatesdp_blacklist_user" "test" { reason = "x" }`,
			ExpectError: regexp.MustCompile(`user_distinguished_name`),
		}},
	})
}

// TestReplicationSource: registering is a POST of the token to the
// /replication-source singleton (v24+); destroy unregisters, tolerating an
// already disabled source.
func TestReplicationSource(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		step := resource.TestStep{Config: providerConfig(srv, v.APIVersion) + `
resource "appgatesdp_replication_source" "test" {
  token = "registration-token"
}
`}
		if !supportedIn(v, "v24", "v25") {
			step.ExpectError = regexp.MustCompile(`appgatesdp_replication_source is not available in SDP Admin API ` + v.APIVersion)
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: protoV6Factories, Steps: []resource.TestStep{step}})
			return
		}
		step.Check = resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("appgatesdp_replication_source.test", "id", "replication_source"),
			resource.TestCheckResourceAttr("appgatesdp_replication_source.test", "token", "registration-token"),
			func(*terraform.State) error {
				if body := lastRequest(srv, http.MethodPost, "/replication-source"); body["token"] != "registration-token" {
					return fmt.Errorf("POST body %v", body)
				}
				// The source is disabled outside Terraform before destroy.
				srv.Fail(http.MethodDelete, "/replication-source", http.StatusPreconditionFailed, `{"id":"precondition failed","message":"Replication source is not enabled."}`, 1)
				return nil
			},
		)
		resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: protoV6Factories, Steps: []resource.TestStep{step}})
	})
}

// TestControllerActivation enables the controller on an activated
// appliance, waits for controller_ready and disables it on destroy.
func TestControllerActivation(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		srv.PendingStatusPolls = 2
		id := srv.Seed("/appliances", seededAppliance("ctrl-2", true))
		cfg := providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_appliance_controller_activation" "test" {
  appliance_id = %q
  controller = {
    enabled = true
  }
  admin_interface = {
    hostname      = "ctrl-2.example.com"
    https_ciphers = ["ECDHE-RSA-AES256-GCM-SHA384"]
  }
}
`, id)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if ctrl := srv.Object("/appliances", id)["controller"].(map[string]any); ctrl["enabled"] != false {
					return fmt.Errorf("destroy must disable the controller, got %v", ctrl)
				}
				if srv.Object("/appliances", id) == nil {
					return fmt.Errorf("the appliance itself must not be deleted")
				}
				return nil
			},
			Steps: []resource.TestStep{
				{
					Config: cfg,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("appgatesdp_appliance_controller_activation.test", "id", id),
						resource.TestCheckResourceAttr("appgatesdp_appliance_controller_activation.test", "controller.enabled", "true"),
						func(*terraform.State) error {
							obj := srv.Object("/appliances", id)
							if obj["controller"].(map[string]any)["enabled"] != true || obj["name"] != "ctrl-2" {
								return fmt.Errorf("PUT must enable the controller and keep the appliance: %v", obj)
							}
							polls := 0
							for _, r := range srv.Requests() {
								if r.Path == "/appliances/status" {
									polls++
								}
							}
							if polls < 3 {
								return fmt.Errorf("expected to wait for controller_ready, polled %d times", polls)
							}
							return nil
						},
					),
				},
				{ResourceName: "appgatesdp_appliance_controller_activation.test", ImportState: true, ImportStateId: id, ImportStateVerify: true},
			},
		})
	})
}

func TestControllerActivationErrors(t *testing.T) {
	srv := newServer(t, "v25")
	inactive := srv.Seed("/appliances", map[string]any{"name": "new", "activated": false})
	active := srv.Seed("/appliances", seededAppliance("slow", true))
	srv.PendingStatusPolls = 1 << 30
	cfg := func(id, extra string) string {
		return providerConfig(srv, "v25") + fmt.Sprintf(`
resource "appgatesdp_appliance_controller_activation" "test" {
  appliance_id = %q
  %s
  controller      = { enabled = true }
  admin_interface = {
    hostname      = "x"
    https_ciphers = ["c"]
  }
}
`, id, extra)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{
			{Config: cfg(inactive, ""), ExpectError: regexp.MustCompile(`(?s)seed\s+appliance\s+"new"\s+first`)},
			{Config: cfg(active, `wait_timeout = "soon"`), ExpectError: regexp.MustCompile(`Invalid wait_timeout`)},
			// Timing out is a warning: the change is applied and kept in state.
			{Config: cfg(active, `wait_timeout = "50ms"`), Check: resource.TestCheckResourceAttr("appgatesdp_appliance_controller_activation.test", "controller.enabled", "true")},
		},
	})
}
