package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateLocalUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                resource "appgate_local_user" "test_local_user" {
                    name                  = "testapiuser"
                    first_name            = "john"
                    last_name             = "doe"
                    password              = "hunter2"
                    email                 = "john.doe@test.com"
                    phone                 = "+1-202-555-0172"
                    failed_login_attempts = 30
                    lock_start            = "2020-04-27T09:51:03Z"
                    tags = [
                      "terraform",
                      "api-created"
                    ]
                }
                data "appgate_local_user" "test" {
                    depends_on = [
                        appgate_local_user.test_local_user,
                    ]
                    local_user_name = "testapiuser"
                }
                `,
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_local_user.test", "local_user_name"),
					resource.TestCheckResourceAttrSet("data.appgate_local_user.test", "local_user_id"),
				),
			},
		},
	})
}
