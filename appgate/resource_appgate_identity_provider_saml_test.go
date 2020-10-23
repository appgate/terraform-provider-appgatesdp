package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSamlIdentityProviderBasic(t *testing.T) {
	resourceName := "appgate_saml_identity_provider.saml_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSamlIdentityProviderBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP"),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.attribute_name", "objectGUID"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.claim_name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.attribute_name", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.attribute_name", "givenName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.claim_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.attribute_name", "sn"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.attribute_name", "mail"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.claim_name", "emails"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.attribute_name", "memberOf"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.claim_name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "6"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.claim_name", "antiVirusIsRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.command", "fileSize"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.path", "/usr/bin/python3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.company.com"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccSamlIdentityProviderImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckSamlIdentityProviderBasic(rName string) string {
	return fmt.Sprintf(`
data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}
data "appgate_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgate_saml_identity_provider" "saml_test_resource" {
  name = "%s"

  admin_provider = true
  ip_pool_v4     = data.appgate_ip_pool.ip_v4_pool.id
  ip_pool_v6     = data.appgate_ip_pool.ip_v6_pool.id
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31",
    "192.100.111.32",
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  redirect_url = "https://saml.company.com"
  issuer       = "http://adfs-test.company.com/adfs/services/trust"
  audience     = "Company Appgate SDP"


  provider_certificate = <<-EOF
-----BEGIN CERTIFICATE-----
MIICZjCCAc+gAwIBAgIUT0AsBLRI7aKjaMTnH1N9J6eS+7EwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDA5MjIxNDQ5MTZaFw0yMTA5
MjIxNDQ5MTZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBAOWp5CnfLvNpjeESzTg/B/1kG1BRdXtM00q59WPj7adZ5gq+
+Hr0mWEQ5GldgmXRE3HsXfv7hiq4RwX9h+qtRinwhSvtLquM54/Fpw+TYZl5N27m
ov8a04qqlo8c3BqXR5Vp+ohPVcXs2I21k5bUTh5XwHj4uiv8uxmKzk42WETbAgMB
AAGjUzBRMB0GA1UdDgQWBBSpc1YN7rgPiBrVPn0roGV+1B4ETDAfBgNVHSMEGDAW
gBSpc1YN7rgPiBrVPn0roGV+1B4ETDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4GBAMgxxBlfgH98ME7Es9xlV3HrurwG1p2gBvrrEACMtFNgtZE1vgck
jmhbc3t+Af9Dv9KBkaI6ZDl16uiptdpAv59wLgbVFgEPUJboRjhIaw5mPcMCeSDE
eIE/AV/qHWNEiLIMP5JO2FUbjpDCYtHkCOFDmv01e6rs86L3MQ8zF76T
-----END CERTIFICATE-----
EOF

  block_local_dns_requests = true
  on_boarding_two_factor {
    mfa_provider_id       = data.appgate_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  tags = [
    "terraform",

  ]
  claim_mappings {
    attribute_name = "objectGUID"
    claim_name     = "userId"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sAMAccountName"
    claim_name     = "username"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "givenName"
    claim_name     = "firstName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sn"
    claim_name     = "lastName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "mail"
    claim_name     = "emails"
    encrypted      = false
    list           = true
  }
  claim_mappings {
    attribute_name = "memberOf"
    claim_name     = "groups"
    encrypted      = false
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

func testAccCheckSamlIdentityProviderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SamlIdentityProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching saml identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSamlIdentityProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_saml_identity_provider" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SamlIdentityProvidersApi

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("saml identity provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccSamlIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
