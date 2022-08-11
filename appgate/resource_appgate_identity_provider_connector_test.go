package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConnectorIdentityProviderBasic(t *testing.T) {
	resourceName := "appgatesdp_connector_identity_provider.connector_test_resource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConnectorIdentityProviderBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", identityProviderConnector),
					resource.TestCheckResourceAttrSet(resourceName, "claim_mappings.#"),
					resource.TestCheckResourceAttr(resourceName, "name", builtinProviderConnector),
					resource.TestCheckResourceAttr(resourceName, "notes", "Built-in Identity Provider on local database."),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "builtin"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccConnectorIdentityProviderImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckConnectorIdentityProviderBasic() string {
	return `
data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

resource "appgatesdp_connector_identity_provider" "connector_test_resource" {
  notes      = "Built-in Identity Provider on local database."
  ip_pool_v4 = data.appgatesdp_ip_pool.ip_v4_pool.id
  ip_pool_v6 = data.appgatesdp_ip_pool.ip_v6_pool.id
  tags = [
    "builtin",
  ]
}
`
}

func testAccConnectorIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
