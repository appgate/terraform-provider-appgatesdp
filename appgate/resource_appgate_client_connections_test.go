package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccClientConnectionsBasic(t *testing.T) {
	resourceName := "appgate_client_connections.test_example_client_connections"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckClientConnectionsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientConnectionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "spa_mode"),
					resource.TestCheckResourceAttr(resourceName, "profiles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profiles.0.identity_provider_name", "local"),
					resource.TestCheckResourceAttr(resourceName, "profiles.0.name", "Company Test"),
					resource.TestCheckResourceAttr(resourceName, "profiles.0.spa_key_name", "test_key"),
					resource.TestCheckResourceAttrSet(resourceName, "profiles.0.url"),
					resource.TestCheckResourceAttr(resourceName, "spa_mode", "TCP"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccClientConnectionsImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckClientConnectionsBasic() string {
	return fmt.Sprintf(`
    resource "appgate_client_connections" "test_example_client_connections" {
        spa_mode = "TCP"
        profiles {
          name                   = "Company Test"
          spa_key_name           = "test_key"
          identity_provider_name = "local"
        }
      }
`)
}

func testAccCheckClientConnectionsExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ClientConnectionsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.ClientConnectionsGet(context.Background()).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching ClientConnections with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccClientConnectionsImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
