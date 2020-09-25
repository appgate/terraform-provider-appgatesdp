package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccEntitlementBasicPing(t *testing.T) {
	resourceName := "appgate_entitlement.test_item"
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

					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.0", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.1", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.2", "hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.3", "dns://hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.hosts.4", "aws://security-group:accounting"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.ports.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.subtype", "icmp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.3805508908.types.0", "0-16"),

					resource.TestCheckResourceAttr(resourceName, "app_shortcut.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcut.1872821293.color_code", "5"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcut.1872821293.name", "ping"),
					resource.TestCheckResourceAttr(resourceName, "app_shortcut.1872821293.url", "https://www.google.com"),
					resource.TestCheckResourceAttr(resourceName, "condition_logic", "and"),

					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2876187004", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
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
		if rs.Type != "appgate_entitlement" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.EntitlementsApi

		_, _, err := api.EntitlementsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Entitlement still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckEntitlementBasicPing(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}
data "appgate_condition" "always" {
  condition_name = "Always"
}
resource "appgate_entitlement" "test_item" {
  name        = "%s"
  site = data.appgate_site.default_site.id
    conditions = [
      data.appgate_condition.always.id
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

  app_shortcut {
    name       = "ping"
    url        = "https://www.google.com"
    color_code = 5
  }
}
`, rName)
}

func testAccCheckEntitlementExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.EntitlementsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.EntitlementsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
