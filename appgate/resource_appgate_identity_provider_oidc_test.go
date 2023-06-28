package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOidcIdentityProviderBasic(t *testing.T) {
	resourceName := "appgatesdp_oidc_identity_provider.oidc_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOidcIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor62AndAbove(t)
				},
				Config: testAccCheckOidcIdentityProviderBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOidcIdentityProviderExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "type", "Oidc"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://example.com/oidc/issuer"),
					resource.TestCheckResourceAttr(resourceName, "audience", "oidc_test_audience"),
					resource.TestCheckResourceAttr(resourceName, "scope", "oidc_test_scope"),
					resource.TestCheckResourceAttr(resourceName, "google.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "google.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "google.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "google.0.client_secret", "oidc_test_client_secret"),
					resource.TestCheckResourceAttr(resourceName, "google.0.refresh_token", "true"),

					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.attribute_name", "givenName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.claim_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.attribute_name", "mail"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.claim_name", "emails"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.attribute_name", "memberOf"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.claim_name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.attribute_name", "objectGUID"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.claim_name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.attribute_name", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.attribute_name", "sn"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "device_limit_per_user", "77"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.claim_name", "antiVirusIsRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.command", "fileSize"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.path", "/usr/bin/python3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccOidcIdentityProviderImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"google.0.client_secret"},
			},
		},
	})
}

func testAccCheckOidcIdentityProviderBasic(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_oidc_identity_provider" "oidc_test_resource" {
  issuer = "https://example.com/oidc/issuer"
  audience = "oidc_test_audience"
  scope = "oidc_test_scope"
  google { 
    enabled = true
    client_secret = "oidc_test_client_secret"
    refresh_token = true
  }

  name = "%s"
  admin_provider = true
  ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
  ip_pool_v6     = data.appgatesdp_ip_pool.ip_v6_pool.id
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31"
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  block_local_dns_requests = true
  device_limit_per_user = 77
  on_boarding_two_factor {
    mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
    message               = "welcome"
  }
  tags = [
    "terraform",
    "api-created"
  ]
  claim_mappings {
    attribute_name = "objectGUID"
    claim_name     = "userId"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sAMAccountName"
    claim_name     = "username"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "givenName"
    claim_name     = "firstName"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sn"
    claim_name     = "lastName"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "mail"
    claim_name     = "emails"
    encrypt      = false
    list           = true
  }
  claim_mappings {
    attribute_name = "memberOf"
    claim_name     = "groups"
    encrypt      = false
    list           = true
  }

  on_demand_claim_mappings {
    command    = "fileSize"
    claim_name = "antiVirusIsRunning"
    parameters {
      path = "/usr/bin/python3"
    }
    platform = "desktop.windows.all"
  }
}
`, rName)
}

func testAccCheckOidcIdentityProviderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.OidcIdentityProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching Oidc identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckOidcIdentityProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_oidc_identity_provider" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.OidcIdentityProvidersApi

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("oidc identity provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccOidcIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
