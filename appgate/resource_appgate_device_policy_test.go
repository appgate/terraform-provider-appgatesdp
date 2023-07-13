package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyDeviceBasic(t *testing.T) {
	resourceName := "appgatesdp_device_policy.test_device_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyDeviceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "11"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Hide"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Hide"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "terraform policy notes"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", "http://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.0", "developer"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "ringfence_rules.*", "data.appgatesdp_ringfence_rule.default_ringfence_rule", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", "aa"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "Device"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyDeviceBasicUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "11"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Hide"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "terraform policy notes"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", "http://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", "aa"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "Device"),
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

func testAccCheckPolicyDeviceBasic(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_ringfence_rule" "default_ringfence_rule" {
	ringfence_rule_name = "Block-in"
}

resource "appgatesdp_device_policy" "test_device_policy" {
    name  = "%s"
    notes = "terraform policy notes"
    tags = [
        "terraform",
        "api-created"
    ]

    ringfence_rule_links = [
        "developer"
    ]
	ringfence_rules = [
		data.appgatesdp_ringfence_rule.default_ringfence_rule.id
	]
    tamper_proofing = true
    proxy_auto_config {
        enabled = true
        url     = "http://foo.com"
        persist = false
    }
    trusted_network_check {
        enabled    = true
        dns_suffix = "aa"
    }
	client_settings {
		enabled           = true
		entitlements_list = "Hide"
		quit              = "Hide"
	}
}

`, rName)
}
func testAccCheckPolicyDeviceBasicUpdated(rName string) string {
	return fmt.Sprintf(`

resource "appgatesdp_device_policy" "test_device_policy" {
    name  = "%s"
    notes = "terraform policy notes"
    tags = [
        "api-created"
    ]

    tamper_proofing = true
    proxy_auto_config {
        enabled = true
        url     = "http://foo.com"
        persist = false
    }
    trusted_network_check {
        enabled    = true
        dns_suffix = "aa"
    }
	client_settings {
		enabled           = true
		entitlements_list = "Hide"
		quit              = "Show"
	}
}

`, rName)
}
