package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var resourceName = "appgate_site.test_site"

func TestAccSiteBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSite(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", "The test site"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "tst"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),

					resource.TestCheckResourceAttr(
						resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.1941342072", "developer"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.2876187004", "api-created"),

					resource.TestCheckResourceAttr(
						resourceName, "default_gateway.1013162454.enabled_v4", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "default_gateway.1013162454.enabled_v6", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "default_gateway.1013162454.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "network_subnets.2657174977", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(
						resourceName, "vpn.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "name_resolution.#", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "ip_pool_mappings.#", "0"),
				),
			},
		},
	})
}

func testAccSiteImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %#v", expectedStates, len(s), s)
		}
		return nil
	}
}

func testAccCheckSiteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_site" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SitesApi

		_, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Site still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckSite() string {
	return fmt.Sprintf(`
resource "appgate_site" "test_site" {
  name       = "The test site"
  short_name = "tst"
  tags = [
    "developer",
    "api-created"
  ]

  notes = "This object has been created for test purposes."

  network_subnets = [
    "10.0.0.0/16"
  ]
}
`)
}

func testAccCheckSiteExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SitesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
