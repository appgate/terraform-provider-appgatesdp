package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSamlIdentityProviderBasic(t *testing.T) {
	resourceName := "appgatesdp_saml_identity_provider.saml_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.GreaterThanOrEqual(Appliance55Version) {
						t.Skip("Test only for 5.4 and below, on_boarding_two_factor.0.device_limit_per_user updated behaviour in > 5.5")
					}
				},
				Config: testAccCheckSamlIdentityProviderBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.company.com"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"),
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "6"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
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
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccSamlIdentityProviderImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSamlIdentityProviderUpdates(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName), //Name cannot be changed
					resource.TestCheckResourceAttr(resourceName, "notes", "Test note change"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.21"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.3", "172.17.18.20"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "update.company.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.1", "test.company.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.update.company.com"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.update.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP - Update"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"), //Write custom function to check for content of certificate?
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "false"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "5"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "change"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "change"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraformm"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "7"),
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
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.6.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.6.attribute_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.6.claim_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.6.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.6.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.claim_name", "anotherOne"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.command", "serviceRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.name", "test2.exe"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.path", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.claim_name", "antiVirusIsRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.command", "fileSize"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.path", "/usr/bin/python3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.claim_name", "testing"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.command", "serviceRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.parameters.0.name", "test.exe"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.parameters.0.path", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.2.platform", "desktop.windows.all"),
				),
			},
			{
				Config: testAccCheckSamlIdentityProviderClaimMoveAndDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Test note change"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "172.17.18.20"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "test.company.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.update.company.com"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.update.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP - Update"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"), //Write custom function to check for content of certificate?
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "false"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "5"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "change"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "5"),
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
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.attribute_name", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.attribute_name", "sn"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.encrypt", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.attribute_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.claim_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.encrypt", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.claim_name", "anotherOne"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.command", "serviceRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.name", "test2.exe"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.parameters.0.path", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.0.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.claim_name", "testing"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.command", "serviceRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.name", "test.exe"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.parameters.0.path", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.1.platform", "desktop.windows.all"),
				),
			},
			{
				// Make sure moving around claims results in an empty plan
				Config: testAccCheckSamlIdentityProviderClaimMoves(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCheckSamlIdentityProviderBasic(rName string) string {
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
	resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
	name = "%s"

	admin_provider = true
	ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
	ip_pool_v6     = data.appgatesdp_ip_pool.ip_v6_pool.id
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
		mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
		device_limit_per_user = 6
		message               = "welcome"
	}
	tags = [
		"terraform",

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

func testAccCheckSamlIdentityProviderUpdates(rName string) string {
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
	resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
	name = "%s"
	notes = "Test note change"
	admin_provider = false
	ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
	dns_servers = [
		"172.17.18.21",
		"192.100.111.31",
		"192.100.111.32",
		"172.17.18.20"
	]
	dns_search_domains = [
		"update.company.com",
		"test.company.com"
	]
	redirect_url = "https://saml.update.company.com"
	issuer       = "http://adfs-test.update.company.com/adfs/services/trust"
	audience     = "Company Appgate SDP - Update"


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

	block_local_dns_requests = false
	inactivity_timeout_minutes = 5

	on_boarding_two_factor {
		mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
		device_limit_per_user = 4
		message               = "change"
	}
	tags = [
		"terraformm",
		"change"

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
		claim_mappings {
		attribute_name = "test"
		claim_name     = "test"
		encrypt      = false
		list           = false
	}
	on_demand_claim_mappings {
		command    = "fileSize"
		claim_name = "antiVirusIsRunning"
		parameters {
		path = "/usr/bin/python3"
		}
		platform = "desktop.windows.all"
	}
	on_demand_claim_mappings {
		command    = "serviceRunning"
		claim_name = "testing"
		parameters {
		name = "test.exe"
		}
		platform = "desktop.windows.all"
	}
	on_demand_claim_mappings {
		command    = "serviceRunning"
		claim_name = "anotherOne"
		parameters {
		name = "test2.exe"
		}
		platform = "desktop.windows.all"
	}
	}`, rName)
}

func testAccCheckSamlIdentityProviderClaimMoveAndDelete(rName string) string {
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
  resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
    name = "%s"
    notes = "Test note change"
    admin_provider = false
    ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
    dns_servers = [
      "192.100.111.32",
      "172.17.18.20"
    ]
    dns_search_domains = [
      "test.company.com"
    ]
    redirect_url = "https://saml.update.company.com"
    issuer       = "http://adfs-test.update.company.com/adfs/services/trust"
    audience     = "Company Appgate SDP - Update"


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

    block_local_dns_requests = false
    inactivity_timeout_minutes = 5
    on_boarding_two_factor {
      mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
      device_limit_per_user = 4
      message               = "change"
    }

    claim_mappings {
      attribute_name = "sn"
      claim_name     = "lastName"
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
      attribute_name = "mail"
      claim_name     = "emails"
      encrypt      = false
      list           = true
    }
      claim_mappings {
      attribute_name = "test"
      claim_name     = "test"
      encrypt      = true
      list           = false
    }
    claim_mappings {
      attribute_name = "givenName"
      claim_name     = "firstName"
      encrypt      = false
      list           = false
    }
    on_demand_claim_mappings {
      command    = "serviceRunning"
      claim_name = "anotherOne"
      parameters {
        name = "test2.exe"
      }
      platform = "desktop.windows.all"
    }
    on_demand_claim_mappings {
      command    = "serviceRunning"
      claim_name = "testing"
      parameters {
        name = "test.exe"
      }
      platform = "desktop.windows.all"
    }

  }`, rName)
}

func testAccCheckSamlIdentityProviderClaimMoves(rName string) string {
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
  resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
    name = "%s"
    notes = "Test note change"
    admin_provider = false
    ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
    dns_servers = [
      "192.100.111.32",
      "172.17.18.20"
    ]
    dns_search_domains = [
      "test.company.com"
    ]
    redirect_url = "https://saml.update.company.com"
    issuer       = "http://adfs-test.update.company.com/adfs/services/trust"
    audience     = "Company Appgate SDP - Update"


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

    block_local_dns_requests = false
    inactivity_timeout_minutes = 5
    on_boarding_two_factor {
      mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
      device_limit_per_user = 4
      message               = "change"
    }
	claim_mappings {
		attribute_name = "mail"
		claim_name     = "emails"
		encrypt      = false
		list           = true
	  }
	claim_mappings {
		attribute_name = "test"
		claim_name     = "test"
		encrypt      = true
		list           = false
	}
    claim_mappings {
      attribute_name = "sn"
      claim_name     = "lastName"
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
    on_demand_claim_mappings {
      command    = "serviceRunning"
      claim_name = "testing"
      parameters {
        name = "test.exe"
      }
      platform = "desktop.windows.all"
    }
    on_demand_claim_mappings {
		command    = "serviceRunning"
		claim_name = "anotherOne"
		parameters {
		  name = "test2.exe"
		}
		platform = "desktop.windows.all"
	  }
  }`, rName)
}

func testAccCheckSamlIdentityProviderBasic55OrGreater(rName string) string {
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
	resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
	name = "%s"

	admin_provider = true
	ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
	ip_pool_v6     = data.appgatesdp_ip_pool.ip_v6_pool.id
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
	device_limit_per_user = 44
	on_boarding_two_factor {
		mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
		message               = "welcome"
	}
	tags = [
		"terraform",

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
func TestAccSamlIdentityProviderBasic55OrGreater(t *testing.T) {
	resourceName := "appgatesdp_saml_identity_provider.saml_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.LessThan(Appliance55Version) {
						t.Skip("Test only for 5.5 and above, on_boarding_two_factor.0.device_limit_per_user updated behaviour in > 5.5")
					}
				},
				Config: testAccCheckSamlIdentityProviderBasic55OrGreater(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.company.com"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"),
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "device_limit_per_user", "44"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
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

// TestAccSamlIdentityProviderUserScripts55OrGreater tests https://github.com/appgate/terraform-provider-appgatesdp/issues/246
func TestAccSamlIdentityProviderUserScripts55OrGreater(t *testing.T) {
	resourceName := "appgatesdp_saml_identity_provider.saml_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.LessThan(Appliance55Version) {
						t.Skip("Test only for 5.5 and above, on_boarding_two_factor.0.device_limit_per_user updated behaviour in > 5.5")
					}
				},
				Config: testAccCheckSamlIdentityProviderUserScripts55OrGreater(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.2", "192.100.111.32"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "https://saml.company.com"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "http://adfs-test.company.com/adfs/services/trust"),
					resource.TestCheckResourceAttr(resourceName, "audience", "Company Appgate SDP"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_certificate"),
					resource.TestCheckResourceAttr(resourceName, "decryption_key", ""),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "device_limit_per_user", "44"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Saml"),
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
					resource.TestCheckResourceAttr(resourceName, "user_scripts.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "user_scripts.0", "appgatesdp_user_claim_script.custom_script", "id"),
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

func testAccCheckSamlIdentityProviderUserScripts55OrGreater(rName string) string {
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

resource "appgatesdp_user_claim_script" "custom_script" {
	name       = "user claim script"
	notes      = "This object has been created for test purposes."
	expression = "return {'posture': 25};"
	tags = [
		"developer",
		"api-created"
	]
}

resource "appgatesdp_saml_identity_provider" "saml_test_resource" {
	depends_on = [
		appgatesdp_user_claim_script.custom_script
	]
	name           = "%s"
	admin_provider = true
	ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
	ip_pool_v6     = data.appgatesdp_ip_pool.ip_v6_pool.id
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

	user_scripts = [
		appgatesdp_user_claim_script.custom_script.id
	]

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
	device_limit_per_user    = 44
	on_boarding_two_factor {
		mfa_provider_id = data.appgatesdp_mfa_provider.fido.id
		message         = "welcome"
	}
	tags = [
		"terraform",

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

func testAccCheckSamlIdentityProviderExists(resource string) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.SamlIdentityProvidersApi
		rs, ok := state.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching saml identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSamlIdentityProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_saml_identity_provider" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.SamlIdentityProvidersApi

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
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
