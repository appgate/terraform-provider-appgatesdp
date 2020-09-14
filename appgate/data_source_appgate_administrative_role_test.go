package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateAdministrativeRoleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                resource "appgate_administrative_role" "test_administrative_role" {
					name  = "tf-acceptance-test-admin-role"
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
				data "appgate_administrative_role" "test" {
                    depends_on = [
                        appgate_administrative_role.test_administrative_role,
                    ]
                    administrative_role_name = "tf-acceptance-test-admin-role"
                }
                `,
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_administrative_role.test", "administrative_role_name"),
					resource.TestCheckResourceAttrSet("data.appgate_administrative_role.test", "administrative_role_id"),
				),
			},
		},
	})
}
