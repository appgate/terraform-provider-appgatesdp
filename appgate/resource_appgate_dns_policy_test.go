package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyDnsBasic(t *testing.T) {
	resourceName := "appgatesdp_dns_policy.test_dns_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":                rName,
		"entitlment_one_name": RandStringFromCharSet(10, CharSetAlphaNum),
		"entitlment_two_name": RandStringFromCharSet(10, CharSetAlphaNum),
		"new_name":            rName + "NEW",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyDnsBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.domain", "appgate.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "Dns"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "entitlements.*", "appgatesdp_entitlement.one", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "entitlements.*", "appgatesdp_entitlement.two", "id"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyDnsBasicUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.domain", "google.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.0", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.1", "3.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", context["new_name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "Dns"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckPolicyDnsBasic(context map[string]interface{}) string {
	return Nprintf(`
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "one" {
	app_shortcut_scripts = []
	condition_logic      = "and"
	conditions = [
		data.appgatesdp_condition.always.id,
	]
	name = "%{entitlment_one_name}"
	site = data.appgatesdp_site.default_site.id

	actions {
		action = "allow"
		hosts = [
		"google.se",
		]
		ports = [
		"443",
		"80",
		]
		subtype = "tcp_up"

	}
}
resource "appgatesdp_entitlement" "two" {
	app_shortcut_scripts = []
	condition_logic      = "and"
	conditions = [
		data.appgatesdp_condition.always.id,
	]
	name = "%{entitlment_two_name}"
	site = data.appgatesdp_site.default_site.id
	tags = []

	actions {
		action = "allow"
		hosts = [
		"appgate.com",
		]
		ports = [
		"443",
		"80",
		]
		subtype = "tcp_up"
	}
}

resource "appgatesdp_dns_policy" "test_dns_policy" {
	name = "%{name}"
	dns_settings {
		domain  = "appgate.com"
		servers = ["8.8.8.8", "1.1.1.1"]
	}
	entitlements = [
		appgatesdp_entitlement.one.id,
		appgatesdp_entitlement.two.id,
	]
}
`, context)
}

func testAccCheckPolicyDnsBasicUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_dns_policy" "test_dns_policy" {
	name = "%{new_name}"
	dns_settings {
		domain  = "google.com"
		servers = ["2.2.2.2", "3.3.3.3"]
	}
}
`, context)
}
