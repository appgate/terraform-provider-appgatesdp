package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateAdministrativeRoleDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgate_administrative_role.test_administrative_role_ds"
	resourceName := "appgate_administrative_role.test_administrative_role"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: administrativeRoleDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "administrative_role_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "administrative_role_id", resourceName, "id"),
				),
			},
		},
	})
}

func administrativeRoleDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "appgate_administrative_role" "test_administrative_role" {
    name  = "%s"
    notes = "hello world"
    tags = [
        "terraform"
    ]
    privileges {
        type         = "Create"
        target       = "Entitlement"
        default_tags = ["api-created"]
    }
}
data "appgate_administrative_role" "test_administrative_role_ds" {
    administrative_role_id = appgate_administrative_role.test_administrative_role.id
}`, rName)
}
