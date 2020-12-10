package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateAdministrativeRoleDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: administrativeRoleDataSourceConfig(rName),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				//ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_administrative_role.test_administrative_role_ds", "administrative_role_name"),
					resource.TestCheckResourceAttrSet("data.appgate_administrative_role.test_administrative_role_ds", "administrative_role_id"),
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
    depends_on = [
        appgate_administrative_role.test_administrative_role,
    ]
    administrative_role_name = "%s"
}`, rName, rName)
}
