package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyAdminBasic(t *testing.T) {
	resourceName := "appgatesdp_admin_policy.test_admin_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":           rName,
		"admin_name_one": RandStringFromCharSet(10, CharSetAlphaNum),
		"admin_name_two": RandStringFromCharSet(10, CharSetAlphaNum),
		"new_name":       rName + "NEW",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyAdminBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "bb"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "cc"),
					resource.TestCheckResourceAttr(resourceName, "type", "Admin"),
					resource.TestCheckResourceAttr(resourceName, "administrative_roles.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "administrative_roles.*", "appgatesdp_administrative_role.first", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "administrative_roles.*", "appgatesdp_administrative_role.second", "id"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyAdminBasicUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", context["new_name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "administrative_roles.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "hello"),
					resource.TestCheckResourceAttr(resourceName, "type", "Admin"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckPolicyAdminBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_administrative_role" "first" {
	name  ="%{admin_name_one}"
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
resource "appgatesdp_administrative_role" "second" {
	name  ="%{admin_name_two}"
	notes = "hello world"
	tags = [
		"terraform"
	]
	privileges {
		type   = "View"
		target = "All"
	}
}

resource "appgatesdp_admin_policy" "test_admin_policy" {
	name = "%{name}"
	tags = ["aa", "bb", "cc"]
	administrative_roles = [
		appgatesdp_administrative_role.first.id,
		appgatesdp_administrative_role.second.id,
	]
}
`, context)
}

func testAccCheckPolicyAdminBasicUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_admin_policy" "test_admin_policy" {
	name = "%{new_name}"
	tags = ["hello"]
}
`, context)
}
