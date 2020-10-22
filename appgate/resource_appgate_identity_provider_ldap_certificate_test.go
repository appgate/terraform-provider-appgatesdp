package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccLdapCertificateIdentityProvidervBasic(t *testing.T) {
	resourceName := "appgate_ldap_certificate_identity_provider.ldap_cert_test_resource"
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
data "appgate_ip_pool" "ip_four_pool" {
  ip_pool_name = "default pool v4"
}

data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgate_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgate_ldap_certificate_identity_provider" "ldap_cert_test_resource" {
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
  default                    = false
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgate_ip_pool.ip_four_pool.id
  ip_pool_v6                 = data.appgate_ip_pool.ip_v6_pool.id
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
    mfa_provider_id       = data.appgate_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  certificate_user_attribute = "blabla"
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
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.LdapCertificateIdentityProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching ldap identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLdapCertificateIdentityProvidervDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_ldap_certificate_identity_provider" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.LdapCertificateIdentityProvidersApi

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
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
