package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAdminMfaSettingsBasic(t *testing.T) {
	resourceName := "appgate_admin_mfa_settings.test_example_mfa_settings"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminMfaSettingsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminMfaSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "admin_mfa_settings"),
					resource.TestCheckResourceAttr(resourceName, "exempted_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "exempted_users.0", "CN=JohnDoe,OU=local"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccAdminMfaSettingsImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckAdminMfaSettingsBasic() string {
	return fmt.Sprintf(`
resource "appgate_admin_mfa_settings" "test_example_mfa_settings" {
    exempted_users = [
        "CN=JohnDoe,OU=local"
    ]
}
`)
}

func testAccCheckAdminMfaSettingsExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.MFAForAdminsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.AdminMfaSettingsGet(context.Background()).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching AdminMfaSettings with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccAdminMfaSettingsImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
