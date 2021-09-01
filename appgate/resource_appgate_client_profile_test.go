package appgate

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccClientProfileBasic(t *testing.T) {
	resourceName := "appgatesdp_client_profile.test_client_profile"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckClientProfileBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClientProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "hello"),
					resource.TestCheckResourceAttr(resourceName, "spa_key_name", "world"),
					resource.TestCheckResourceAttr(resourceName, "identity_provider_name", "local"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccClientProfileImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckClientProfileBasic() string {
	return `
resource "appgatesdp_client_profile" "test_client_profile" {
	name                   = "hello"
	spa_key_name           = "world"
	identity_provider_name = "local"
}
`
}

func testAccCheckClientProfileExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.ClientConnectionsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No Record ID is set")
		}

		clientConnections, _, err := api.ClientConnectionsGet(context.Background()).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching ClientConnections with resource %s. %s", resource, err)
		}

		for _, profile := range clientConnections.GetProfiles() {
			if strings.EqualFold(profile.GetName(), id) && profile.GetName() == id {
				return nil
			}

		}
		return fmt.Errorf("Could not find client connection.profile with name %s", id)
	}
}

func testAccClientProfileImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
