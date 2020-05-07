package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccEntitlementBasicPing(t *testing.T) {
	resourceName := "appgate_entitlement.test_item"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementBasicPing(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleItemExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", "ping"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.0", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.1", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.2", "hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.3", "dns://hostname.company.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.hosts.4", "aws://security-group:accounting"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.ports.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.subtype", "icmp_up"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4206576320.types.0", "0-16"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// ImportStateVerify: true,
				ImportStateCheck: testAccEntitlementImportStateCheckFunc(1),
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
			return fmt.Errorf("Entitlment still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckEntitlementBasicPing() string {
	// TODO: conditions need to be dynamic, data attribute.
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}
data "appgate_condition" "always" {
  condition_name = "Always"
}
resource "appgate_entitlement" "test_item" {
  name        = "ping"
  site = data.appgate_site.default_site.id
    conditions = [
      data.appgate_condition.always.id
  ]
  actions {
    subtype = "icmp_up"
    action  = "allow"
    types = ["0-16"]
    hosts = [
      "10.0.0.1",
      "10.0.0.0/24",
      "hostname.company.com",
      "dns://hostname.company.com",
      "aws://security-group:accounting"
    ]
  }
}
`)
}

func testAccCheckExampleItemExists(resource string) resource.TestCheckFunc {
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
