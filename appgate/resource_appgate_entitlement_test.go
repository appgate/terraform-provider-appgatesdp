package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEntitlementBasicPing(t *testing.T) {
	resourceName := "appgatesdp_entitlement.test_ping_item"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementBasicPing(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.1", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.2", "aws://security-group:accounting"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.3", "dns://hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.4", "hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "icmp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.0", "0-16"),

					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),

					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccEntitlementImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %#v", expectedStates, len(s), s)
		}
		return nil
	}
}

func testAccCheckItemDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_entitlement" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.EntitlementsApi

		if _, _, err := api.EntitlementsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("Entitlement still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckEntitlementBasicPing(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	   site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
  condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_ping_item" {
  name        = "%s"
  site = data.appgatesdp_site.default_site.id
    conditions = [
      data.appgatesdp_condition.always.id
  ]

  tags = [
    "terraform",
    "api-created"
  ]
  disabled = true

  condition_logic = "and"

  actions {
    subtype = "icmp_up"
    action  = "allow"
    # https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
    types = ["0-16"]
    hosts = [
      "10.0.0.1",
      "10.0.0.0/24",
      "hostname.company.com",
      "dns://hostname.company.com",
      "aws://security-group:accounting"
    ]
  }

  app_shortcuts {
    name       = "%s"
    url        = "https://www.google.com"
    color_code = 5
  }
}
`, rName, rName)
}

func testAccCheckEntitlementExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.EntitlementsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.EntitlementsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func TestAccEntitlementBasicWithMonitor(t *testing.T) {
	resourceName := "appgatesdp_entitlement.monitor_entitlement"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementWithMonitor(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					// default value, computed from the controller, even if we dont set it in the tf
					// config file, we will get a computed value back.
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementWithMonitorUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "22"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementWithMonitor(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	   site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
  condition_name = "Always"
}
resource "appgatesdp_entitlement" "monitor_entitlement" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
	  data.appgatesdp_condition.always.id
	]

	tags = [
	  "terraform",
	  "api-created"
	]
	disabled = true

	condition_logic = "and"
	actions {
	  action  = "allow"
	  subtype = "tcp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	}

	app_shortcuts {
	  name       = "%s"
	  url        = "https://www.google.com"
	  color_code = 5
	}
  }
`, rName, rName)
}

func testAccCheckEntitlementWithMonitorUpdated(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	   site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
  condition_name = "Always"
}
resource "appgatesdp_entitlement" "monitor_entitlement" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
	  data.appgatesdp_condition.always.id
	]

	tags = [
	  "terraform",
	  "api-created"
	]
	disabled = true

	condition_logic = "and"
	actions {
	  action  = "allow"
	  subtype = "tcp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	  monitor {
		enabled = true
		timeout = 22
	  }
	}

	app_shortcuts {
	  name       = "%s"
	  url        = "https://www.google.com"
	  color_code = 5
	}
  }
`, rName, rName)
}

func TestAccEntitlementUpdateActionOrder(t *testing.T) {
	resourceName := "appgatesdp_entitlement.test_action_order_item"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementMultipleActions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),

					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "22"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementWithActionOrderUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementMultipleActions(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_item" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
		monitor {
		enabled = true
		timeout = 22
		}
	}
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 5
	}
}
	`, rName, rName)
}

func testAccCheckEntitlementWithActionOrderUpdated(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_item" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	# Updated the order of the actions
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
		monitor {
			enabled = true
			timeout = 22
		}
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 5
	}
}
	`, rName, rName)
}

func TestAccEntitlementUpdateActionHostOrder(t *testing.T) {
	resourceName := "appgatesdp_entitlement.test_action_order_hosts"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementActionHostSets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "103.15.3.254/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.1", "172.17.3.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.2", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementActionHostSetsUpdatedOrder(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "103.15.3.254/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.1", "172.17.3.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.2", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcut_scripts.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "19"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "conditions.0", "data.appgatesdp_condition.always", "id"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementActionHostSets(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_hosts" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts = [
			"103.15.3.254/32",
			"172.17.3.255/32",
			"192.168.2.255/32",
		]
		ports   = ["53"]
	}
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 5
	}
}
	`, rName, rName)
}

func testAccCheckEntitlementActionHostSetsUpdatedOrder(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_hosts" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts = [
			"192.168.2.255/32",
			"103.15.3.254/32",
			"172.17.3.255/32",
		]
		ports   = ["53"]
	}
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 19
	}
}
	`, rName, rName)
}

func TestAccEntitlementUpdateActionPortsSetOrder(t *testing.T) {
	resourceName := "appgatesdp_entitlement.test_action_order_hosts"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementActionPortsSets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "103.15.3.254/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.1", "172.17.3.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.2", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "21"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.1", "22"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.2", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementActionPortsSetsUpdatedOrder(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "103.15.3.254/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.1", "172.17.3.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.2", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "21"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.1", "22"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.2", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcut_scripts.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "conditions.0", "data.appgatesdp_condition.always", "id"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementActionPortsSets(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_hosts" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts = [
			"103.15.3.254/32",
			"172.17.3.255/32",
			"192.168.2.255/32",
		]
		ports   = ["53", "22", "21"]
	}
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 5
	}
}
	`, rName, rName)
}

func testAccCheckEntitlementActionPortsSetsUpdatedOrder(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_action_order_hosts" {
	name = "%s"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created"
	]
	disabled = true
	condition_logic = "and"
	actions {
		action  = "allow"
		subtype = "tcp_up"
		hosts = [
			"192.168.2.255/32",
			"103.15.3.254/32",
			"172.17.3.255/32",
		]
		ports   = ["21", "22", "53"]
	}
	actions {
		action  = "allow"
		subtype = "udp_up"
		hosts   = ["192.168.2.255/32"]
		ports   = ["53"]
	}
	app_shortcuts {
		name       = "%s"
		url        = "https://www.google.com"
		color_code = 5
	}
}
	`, rName, rName)
}

func TestAccEntitlementRiskSensitivity(t *testing.T) {
	resourceName := "appgatesdp_entitlement.risk"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor6AndAbove(t)
				},
				Config: testAccCheckEntitlementRisk(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "22"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "entitlement_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "risk_sensitivity", "High"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{

				Config: testAccCheckEntitlementRiskUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.0.timeout", "22"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "tcp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.hosts.0", "192.168.2.255/32"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.ports.0", "53"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.subtype", "udp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "entitlement_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "risk_sensitivity", "Low"),
					resource.TestCheckResourceAttrPair(resourceName, "site", "data.appgatesdp_site.default_site", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementRisk(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "risk" {
	name             = "%s"
	site             = data.appgatesdp_site.default_site.id
	conditions       = []
	risk_sensitivity = "high"
	tags = [
	  "terraform",
	  "api-created"
	]
	disabled        = true
	condition_logic = "and"
	actions {
	  action  = "allow"
	  subtype = "tcp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	  monitor {
		enabled = true
		timeout = 22
	  }
	}
	actions {
	  action  = "allow"
	  subtype = "udp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	}
	app_shortcuts {
	  name       = "%s"
	  url        = "https://www.google.com"
	  color_code = 5
	}
}
`, rName, rName)
}

func testAccCheckEntitlementRiskUpdated(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "risk" {
	name             = "%s"
	site             = data.appgatesdp_site.default_site.id
	conditions       = []
	risk_sensitivity = "Low"
	tags = [
	  "terraform",
	  "api-created"
	]
	disabled        = true
	condition_logic = "and"
	actions {
	  action  = "allow"
	  subtype = "tcp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	  monitor {
		enabled = true
		timeout = 22
	  }
	}
	actions {
	  action  = "allow"
	  subtype = "udp_up"
	  hosts   = ["192.168.2.255/32"]
	  ports   = ["53"]
	}
	app_shortcuts {
	  name       = "%s"
	  url        = "https://www.google.com"
	  color_code = 5
	}
}
`, rName, rName)
}

func TestAccEntitlementActionHTTPMethods(t *testing.T) {
	resourceName := "appgatesdp_entitlement.http_methods"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor61AndAbove(t)
				},
				Config: testAccCheckEntitlementHTTPMethods(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.0", "GET"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.1", "HEAD"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.2", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "http_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "18"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementHTTPMethodsUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.0", "POST"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "http_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "18"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckEntitlementHTTPMethodsDeleted(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.methods.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.subtype", "http_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.color_code", "18"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "app_shortcuts.0.url", "http://10.0.5.160"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementHTTPMethods(context map[string]interface{}) string {
	return Nprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "http_methods" {
	condition_logic = "and"

	name = "%{name}"
	site = data.appgatesdp_site.default_site.id
	tags = []
	conditions = []
	actions {
		action = "allow"
		hosts = [
		"http://10.0.5.160",
		]
		methods = [
			"GET",
			"HEAD",
			"PUT",
		]
		ports   = []
		subtype = "http_up"
		types   = []
	}

	app_shortcuts {
		color_code = 18
		name       = "%{name}"
		url        = "http://10.0.5.160"
	}
}

`, context)
}

func testAccCheckEntitlementHTTPMethodsUpdated(context map[string]interface{}) string {
	return Nprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "http_methods" {
	condition_logic = "and"

	name = "%{name}"
	site = data.appgatesdp_site.default_site.id
	tags = []
	conditions = []
	actions {
		action = "allow"
		hosts = [
		"http://10.0.5.160",
		]
		methods = [
			"POST",
		]
		ports   = []
		subtype = "http_up"
		types   = []
	}

	app_shortcuts {
		color_code = 18
		name       = "%{name}"
		url        = "http://10.0.5.160"
	}
}

`, context)
}
func testAccCheckEntitlementHTTPMethodsDeleted(context map[string]interface{}) string {
	return Nprintf(`
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
resource "appgatesdp_entitlement" "http_methods" {
	condition_logic = "and"

	name = "%{name}"
	site = data.appgatesdp_site.default_site.id
	tags = []
	conditions = []
	actions {
		action = "allow"
		hosts = [
		"http://10.0.5.160",
		]
		methods = []
		ports   = []
		subtype = "http_up"
		types   = []
	}

	app_shortcuts {
		color_code = 18
		name       = "%{name}"
		url        = "http://10.0.5.160"
	}
}

`, context)
}
