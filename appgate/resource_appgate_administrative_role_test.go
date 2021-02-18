package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccadministrativeRoleBasic(t *testing.T) {
	resourceName := "appgate_administrative_role.test_administrative_role"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckadministrativeRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckadministrativeRoleBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "hello world"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.0", "cc"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.1", "dd"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Entitlement"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.1", "bb"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.target", "Appliance"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.type", "View"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
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

func testAccCheckadministrativeRoleBasic(rName string) string {
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
        default_tags = ["cc", "dd"]
    }
    privileges {
        type   = "View"
        target = "Appliance"
        scope {
        tags = ["aa", "bb"]
        }
    }
}
`, rName)
}

func testAccCheckadministrativeRoleWithScope(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
  site_name = "Default site"
}
resource "appgate_administrative_role" "administrative_role_with_scope" {
  name = "%s"
  tags = [
    "terraform"
  ]
  privileges {
    type   = "View"
    target = "Site"
    scope {
      ids  = [data.appgate_site.default_site.id]
      tags = ["builtin"]
    }
  }
}
`, rName)
}

func TestAccadministrativeRoleWithScope(t *testing.T) {
	resourceName := "appgate_administrative_role.administrative_role_with_scope"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckadministrativeRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckadministrativeRoleWithScope(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "privileges.0.scope.0.ids.0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.0", "builtin"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Site"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "View"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
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
