package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLdapCertificateIdentityProvidervBasic(t *testing.T) {
	resourceName := "appgatesdp_ldap_certificate_identity_provider.ldap_cert_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapCertificateIdentityProvidervDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLdapCertificateIdentityProvidervBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapCertificateIdentityProvidervExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "admin_distinguished_name", "CN=admin,OU=Users,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "admin_password", "helloworld"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "base_dn", "OU=Users,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "ca_certificates.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "certificate_attribute", "qwerty"),
					resource.TestCheckResourceAttr(resourceName, "certificate_user_attribute", "blabla"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.attribute_name", "givenName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.claim_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.0.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.attribute_name", "mail"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.claim_name", "emails"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.attribute_name", "memberOf"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.claim_name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.attribute_name", "objectGUID"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.claim_name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.attribute_name", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.attribute_name", "sn"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.5.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "dc.ad.company.com"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "28"),
					resource.TestCheckResourceAttr(resourceName, "membership_base_dn", "OU=Groups,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "membership_filter", "(objectCategory=group)"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "object_class", "user"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.always_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.claim_suffix", "onBoarding"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "6"),
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
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "skip_x509_external_checks", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "LdapCertificate"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "sAMAccountName"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccLdapCertificateIdentityProvidervImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

func testAccCheckLdapCertificateIdentityProvidervBasic(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ip_pool" "ip_four_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_ldap_certificate_identity_provider" "ldap_cert_test_resource" {
  name                     = "%s"
  port                     = 389
  admin_distinguished_name = "CN=admin,OU=Users,DC=company,DC=com"
  hostnames                = ["dc.ad.company.com"]
  ssl_enabled              = true
  base_dn                  = "OU=Users,DC=company,DC=com"
  object_class             = "user"
  username_attribute       = "sAMAccountName"
  membership_filter        = "(objectCategory=group)"
  membership_base_dn       = "OU=Groups,DC=company,DC=com"
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgatesdp_ip_pool.ip_four_pool.id
  ip_pool_v6                 = data.appgatesdp_ip_pool.ip_v6_pool.id
  admin_password             = "helloworld"
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31"
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  block_local_dns_requests = true
  on_boarding_two_factor {
    mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  certificate_user_attribute = "blabla"
  certificate_attribute      = "qwerty"
  skip_x509_external_checks  = true
  ca_certificates = [
    <<-EOF
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
  ]
  tags = [
    "terraform",
    "api-created"
  ]
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

func testAccCheckLdapCertificateIdentityProvidervExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LdapCertificateIdentityProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching ldap identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLdapCertificateIdentityProvidervDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_ldap_certificate_identity_provider" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LdapCertificateIdentityProvidersApi

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("ldap identity provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccLdapCertificateIdentityProvidervImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
