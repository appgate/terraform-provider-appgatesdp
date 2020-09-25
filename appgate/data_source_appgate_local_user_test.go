package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateLocalUserDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                resource "appgate_local_user" "new_user_for_ds" {
                    name                  = "%s"
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
                data "appgate_local_user" "testdslu" {
                    depends_on = [
                        appgate_local_user.new_user_for_ds,
                    ]
                    local_user_name = "%s"
                }
                `, rName, rName),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_local_user.testdslu", "local_user_name"),
					resource.TestCheckResourceAttrSet("data.appgate_local_user.testdslu", "local_user_id"),
				),
			},
		},
	})
}
