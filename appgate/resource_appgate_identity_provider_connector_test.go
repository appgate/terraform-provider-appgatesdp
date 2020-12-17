package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConnectorIdentityProviderBasic(t *testing.T) {
	resourceName := "appgate_connector_identity_provider.connector_test_resource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConnectorIdentityProviderBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", identityProviderConnector),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "7"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1018344469.attribute_name", "peerHostname"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1018344469.claim_name", "peerHostname"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1018344469.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1018344469.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.attribute_name", "tags"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.claim_name", "tags"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.attribute_name", "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.claim_name", "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2636647295.attribute_name", "clientName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2636647295.claim_name", "clientName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2636647295.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2636647295.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3447551570.attribute_name", "applianceName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3447551570.claim_name", "applianceName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3447551570.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3447551570.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3602134620.attribute_name", "applianceApiVersion"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3602134620.claim_name", "applianceApiVersion"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3602134620.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3602134620.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4025945261.attribute_name", "hostname"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4025945261.claim_name", "hostname"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4025945261.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4025945261.list", "false"),
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
	return fmt.Sprintf(`
data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

resource "appgate_connector_identity_provider" "connector_test_resource" {
  notes      = "Built-in Identity Provider on local database."
  ip_pool_v4 = data.appgate_ip_pool.ip_v4_pool.id
  ip_pool_v6 = data.appgate_ip_pool.ip_v6_pool.id
  tags = [
    "builtin",
  ]
}
`)
}

func testAccConnectorIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
