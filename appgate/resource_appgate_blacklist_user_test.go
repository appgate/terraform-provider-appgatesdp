package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlacklistUserBasic(t *testing.T) {
	resourceName := "appgatesdp_blacklist_user.test_blacklist_user"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBlacklistUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "CN=TestUser,OU=local"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "local"),
					resource.TestCheckResourceAttr(resourceName, "user_distinguished_name", "CN=TestUser,OU=local"),
					resource.TestCheckResourceAttr(resourceName, "username", "TestUser"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccBlacklistUserImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckBlacklistUserBasic() string {
	return `
    resource "appgatesdp_blacklist_user" "test_blacklist_user" {
        user_distinguished_name = "CN=TestUser,OU=local"
      }
`
}

func testAccBlacklistUserImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
