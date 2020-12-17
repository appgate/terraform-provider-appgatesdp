package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLocalDatabaseIdentityProviderBasic(t *testing.T) {
	resourceName := "appgate_local_database_identity_provider.local_database_test_resource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocalDatabaseIdentityProviderBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "LocalDatabase"),
					resource.TestCheckResourceAttr(resourceName, "min_password_length", "16"),
					resource.TestCheckResourceAttr(resourceName, "user_lockout_threshold", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "admin_provider"),
					resource.TestCheckResourceAttrSet(resourceName, "block_local_dns_requests"),
					resource.TestCheckResourceAttrSet(resourceName, "claim_mappings.#"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1042003740.attribute_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1042003740.claim_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1042003740.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1042003740.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.attribute_name", "tags"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.claim_name", "tags"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1292058217.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1387781900.attribute_name", "phone"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1387781900.claim_name", "phone"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1387781900.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.1387781900.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2348239788.attribute_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2348239788.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2348239788.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2348239788.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.attribute_name", "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.claim_name", "id"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2511141269.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4200917273.attribute_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4200917273.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4200917273.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.4200917273.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.731545266.attribute_name", "email"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.731545266.claim_name", "emails"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.731545266.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.731545266.list", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "default"),
					resource.TestCheckResourceAttrSet(resourceName, "inactivity_timeout_minutes"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v4"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_v6"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "notes"),
					resource.TestCheckResourceAttrSet(resourceName, "on_boarding_two_factor.#"),
					resource.TestCheckResourceAttrSet(resourceName, "on_demand_claim_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.0"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccLocalDatabaseIdentityProviderImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckLocalDatabaseIdentityProviderBasic() string {
	return fmt.Sprintf(`
data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

resource "appgate_local_database_identity_provider" "local_database_test_resource" {
  notes      = "Built-in Identity Provider on local database."
  ip_pool_v4 = data.appgate_ip_pool.ip_v4_pool.id
  ip_pool_v6 = data.appgate_ip_pool.ip_v6_pool.id

  user_lockout_threshold = 10
  min_password_length = 16

  tags = [
    "builtin",
  ]
}
`)
}

func testAccLocalDatabaseIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
