package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateLocalUserDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgatesdp_local_user.testdslu"
	resourceName := "appgatesdp_local_user.new_user_for_ds"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                resource "appgatesdp_local_user" "new_user_for_ds" {
                    name                  = "%s"
                    first_name            = "john"
                    last_name             = "doe"
                    password              = "password_is_hunter2"
                    email                 = "john.doe@test.com"
                    phone                 = "+1-202-555-0172"
                    failed_login_attempts = 30
                    lock_start            = "2020-04-27T09:51:03Z"
                    tags = [
                      "terraform",
                      "api-created"
                    ]
                }
                data "appgatesdp_local_user" "testdslu" {
                    local_user_id = appgatesdp_local_user.new_user_for_ds.local_user_id
                }
                `, rName),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "local_user_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "local_user_id", resourceName, "id"),
				),
			},
		},
	})
}
