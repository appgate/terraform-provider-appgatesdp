package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBlacklistUserBasic(t *testing.T) {
	resourceName := "appgate_blacklist_user.test_blacklist_user"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBlacklistUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "CN=TestUser,OU=ldap"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "ldap"),
					resource.TestCheckResourceAttr(resourceName, "user_distinguished_name", "CN=TestUser,OU=ldap"),
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
	return fmt.Sprintf(`
    resource "appgate_blacklist_user" "test_blacklist_user" {
        user_distinguished_name = "CN=TestUser,OU=ldap"
      }
`)
}

func testAccBlacklistUserImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
