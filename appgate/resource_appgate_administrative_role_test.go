package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccadministrativeRoleBasic(t *testing.T) {
	resourceName := "appgate_administrative_role.test_administrative_role"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckadministrativeRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckadministrativeRoleBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "tf-admin"),
					resource.TestCheckResourceAttr(resourceName, "notes", "hello world"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Entitlement"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccadministrativeRoleImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckadministrativeRoleBasic() string {
	return fmt.Sprintf(`
resource "appgate_administrative_role" "test_administrative_role" {
    name  = "tf-admin"
    notes = "hello world"
    tags = [
        "terraform"
    ]
    privileges {
        type   = "Create"
        target = "Entitlement"
        default_tags = ["api-created"]
    }
}
`)
}

func testAccCheckadministrativeRoleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.AdministrativeRolesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.AdministrativeRolesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching Administrative Role with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckadministrativeRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_administrative_role" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.AdministrativeRolesApi

		_, _, err := api.AdministrativeRolesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Administrative Role still exists, %+v", err)
		}
	}
	return nil
}

func testAccadministrativeRoleImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
