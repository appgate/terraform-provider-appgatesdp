package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLdapIdentityProviderBasic(t *testing.T) {
	resourceName := "appgatesdp_ldap_identity_provider.ldap_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					applianceTestForFiveFiveOrHigher(t)
				},
				Config: testAccCheckLdapIdentityProviderBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_distinguished_name", "CN=admin,OU=Users,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "admin_password", "helloworld"),
					// resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "base_dn", "OU=Users,DC=company,DC=com"),
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
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "dc.ad.company.com"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "28"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
					resource.TestCheckResourceAttr(resourceName, "membership_base_dn", "OU=Groups,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "membership_filter", "(objectCategory=group)"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "object_class", "user"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttrSet(resourceName, "on_boarding_two_factor.0.mfa_provider_id"),
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
					resource.TestCheckResourceAttr(resourceName, "password_warning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.message", "Your password is about to expire, Please change it"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.threshold_days", "13"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Ldap"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "sAMAccountName"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccLdapIdentityProviderImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

func testAccCheckLdapIdentityProviderBasic(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_ldap_identity_provider" "ldap_test_resource" {
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
  password_warning {
    enabled        = true
    threshold_days = 13
    message        = "Your password is about to expire, Please change it"
  }
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgatesdp_ip_pool.ip_v4_pool.id
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
    message               = "welcome"
  }
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

func testAccCheckLdapIdentityProviderBasic55OrGreater(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_ldap_identity_provider" "ldap_test_resource" {
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
  password_warning {
    enabled        = true
    threshold_days = 13
    message        = "Your password is about to expire, Please change it"
  }
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgatesdp_ip_pool.ip_v4_pool.id
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
  device_limit_per_user = 66
  on_boarding_two_factor {
    mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
    message               = "welcome"
  }
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

func TestAccLdapIdentityProviderBasic55OrGreater(t *testing.T) {
	resourceName := "appgatesdp_ldap_identity_provider.ldap_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					applianceConstraintCheck(t, ">= 5.5, < 6.2")
				},
				Config: testAccCheckLdapIdentityProviderBasic55OrGreater(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_distinguished_name", "CN=admin,OU=Users,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "admin_password", "helloworld"),
					// resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "base_dn", "OU=Users,DC=company,DC=com"),
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
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "dc.ad.company.com"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "28"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
					resource.TestCheckResourceAttr(resourceName, "membership_base_dn", "OU=Groups,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "membership_filter", "(objectCategory=group)"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "object_class", "user"),
					resource.TestCheckResourceAttr(resourceName, "device_limit_per_user", "66"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttrSet(resourceName, "on_boarding_two_factor.0.mfa_provider_id"),
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
					resource.TestCheckResourceAttr(resourceName, "password_warning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.message", "Your password is about to expire, Please change it"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.threshold_days", "13"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Ldap"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "sAMAccountName"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccLdapIdentityProviderImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

func testAccCheckLdapIdentityProviderBasic62OrGreater(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_ldap_identity_provider" "ldap_test_resource" {
  name                     = "%s"
  port                     = 389
  admin_distinguished_name = "CN=admin,OU=Users,DC=company,DC=com"
  hostnames                = ["dc.ad.company.com"]
  ssl_enabled              = true
  base_dn                  = "OU=Users,DC=company,DC=com"
  user_filter              = "user"
  username_attribute       = "sAMAccountName"
  membership_filter        = "(objectCategory=group)"
  membership_base_dn       = "OU=Groups,DC=company,DC=com"
  password_warning {
    enabled        = true
    threshold_days = 13
    message        = "Your password is about to expire, Please change it"
  }
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgatesdp_ip_pool.ip_v4_pool.id
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
  device_limit_per_user = 66
  on_boarding_two_factor {
    mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
    message               = "welcome"
  }
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

func TestAccLdapIdentityProviderBasic62OrGreater(t *testing.T) {
	resourceName := "appgatesdp_ldap_identity_provider.ldap_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					applianceConstraintCheck(t, ">= 6.2")
				},
				Config: testAccCheckLdapIdentityProviderBasic62OrGreater(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_distinguished_name", "CN=admin,OU=Users,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "admin_password", "helloworld"),
					// resource.TestCheckResourceAttr(resourceName, "admin_provider", "false"),
					resource.TestCheckResourceAttr(resourceName, "base_dn", "OU=Users,DC=company,DC=com"),
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
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "dc.ad.company.com"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "28"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
					resource.TestCheckResourceAttr(resourceName, "membership_base_dn", "OU=Groups,DC=company,DC=com"),
					resource.TestCheckResourceAttr(resourceName, "membership_filter", "(objectCategory=group)"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "user_filter", "user"),
					resource.TestCheckResourceAttr(resourceName, "device_limit_per_user", "66"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttrSet(resourceName, "on_boarding_two_factor.0.mfa_provider_id"),
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
					resource.TestCheckResourceAttr(resourceName, "password_warning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.message", "Your password is about to expire, Please change it"),
					resource.TestCheckResourceAttr(resourceName, "password_warning.0.threshold_days", "13"),
					resource.TestCheckResourceAttr(resourceName, "port", "389"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Ldap"),
					resource.TestCheckResourceAttr(resourceName, "username_attribute", "sAMAccountName"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccLdapIdentityProviderImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"admin_password"},
			},
		},
	})
}

func testAccCheckLdapIdentityProviderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LdapIdentityProvidersApi

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

func testAccCheckLdapIdentityProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_ldap_identity_provider" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LdapIdentityProvidersApi

		if _, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("ldap identity provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccLdapIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
