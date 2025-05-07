package appgate

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAdministrativeRoleBasic(t *testing.T) {
	resourceName := "appgatesdp_administrative_role.test_administrative_role"
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
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.0", "cc"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Entitlement"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.0", "aa"),
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
			{
				PreConfig: func() {
					testFor6AndAbove(t)
				},
				Config: testAccCheckadministrativeRoleBasicUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.0", "cc"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.1", "dd"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Entitlement"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.target", "Appliance"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.type", "View"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.target", "Ztp"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.type", "View"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.scope.0.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.target", "Appliance"),
					resource.TestCheckResourceAttr(resourceName, "privileges.3.type", "RenewCertificate"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccadministrativeRoleImportStateCheckFunc(1),
			},
			{
				Config:      testAccCheckadministrativeRoleBasicInvalid(rName),
				ExpectError: regexp.MustCompile("Failed to update administrative role privileges scope is not allowed with type View and target Ztp"),
			},
		},
	})
}

func testAccCheckadministrativeRoleBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role" {
    name  = "%s"
    notes = "hello world"
    tags = [
        "terraform"
    ]
    privileges {
        type         = "Create"
        target       = "Entitlement"
        default_tags = ["cc"]
    }
    privileges {
        type   = "View"
        target = "Appliance"
        scope {
        tags = ["aa"]
        }
    }
}
`, rName)
}

func testAccCheckadministrativeRoleBasicUpdated(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role" {
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
        tags = ["aa"]
        }
    }
	privileges {
		type   = "View"
		target = "Ztp"
	  }
	  privileges {
		target = "Appliance"
		type   = "RenewCertificate"

		scope {
		  all = true
		}
	  }
}
`, rName)
}
func testAccCheckadministrativeRoleBasicInvalid(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role" {
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
        tags = ["aa", ]
        }
    }
	privileges {
		type   = "View"
		target = "Ztp"
		// this should not be allowed
		scope {
			all = true
		}
	  }
	  privileges {
		target = "Appliance"
		type   = "RenewCertificate"

		scope {
		  all = true
		}
	  }
}
`, rName)
}

func testAccCheckadministrativeRoleWithScope(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
  site_name = "Default Site"
}
resource "appgatesdp_administrative_role" "administrative_role_with_scope" {
  name = "%s"
  tags = [
    "terraform"
  ]
  privileges {
    type   = "View"
    target = "Site"
    scope {
      ids  = [data.appgatesdp_site.default_site.id]
      tags = ["builtin"]
    }
  }
}
`, rName)
}

func TestAccadministrativeRoleWithScope(t *testing.T) {
	resourceName := "appgatesdp_administrative_role.administrative_role_with_scope"
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
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.AdminRolesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.AdministrativeRolesIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err != nil {
			return fmt.Errorf("error fetching Administrative Role with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckadministrativeRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_administrative_role" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.AdminRolesApi

		if _, _, err := api.AdministrativeRolesIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err == nil {
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

func TestAccadministrativeMultiplePrivilegesValidation(t *testing.T) {
	resourceName := "appgatesdp_administrative_role.test_administrative_role_129"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":   rName,
		"target": "RegisteredDevice", // from >= 5.3 its called RegisteredDevice, prior to 5.3 it was called OnBoardedDevice
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckadministrativeRoleDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
				},
				Config: testAccCheckadministrativeRoleMultiplePrivlegesConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "RegisteredDevice"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Delete"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.target", "RegisteredDevice"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.type", "View"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "bb"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "cc"),
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

func testAccCheckadministrativeRoleMultiplePrivlegesConfig(context map[string]interface{}) string {
	// Test based on https://github.com/appgate/terraform-provider-appgatesdp/issues/129#issuecomment-852211335
	return Nprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role_129" {
	name  = "%{name}"
	tags  = ["aa", "bb", "cc"]

	privileges {
		type   = "View"
		target = "%{target}"

		scope {
		all = true
		}
	}

	privileges {
		type   = "Delete"
		target = "%{target}"

		scope {
		all = true
		}
	}
}
`, context)
}

func TestAccadministrativeRoleWtihAssignFunction(t *testing.T) {
	resourceName := "appgatesdp_administrative_role.test_administrative_role"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckadministrativeRoleDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor6AndAbove(t)
				},
				Config: testAccCheckadministrativeRoleWtihAssignFunction(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "hello"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.0", "cc"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.all", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "Entitlement"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.0", "connector"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.1", "controller"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.2", "gateway"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.functions.3", "logserver"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.scope.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.target", "All"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1.type", "AssignFunction"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.functions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.all", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.scope.0.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.target", "Policy"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2.type", "Create"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccadministrativeRoleImportStateCheckFunc(1),
			},
			{

				Config: testAccCheckadministrativeRoleWtihAssignFunctionUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckadministrativeRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "hello"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.default_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.0", "connector"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.1", "controller"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.2", "gateway"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.functions.3", "logserver"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.scope.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.target", "All"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0.type", "AssignFunction"),
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

func testAccCheckadministrativeRoleWtihAssignFunction(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role" {
	name  = "%{name}"
	notes =  "hello"
	tags = [
	  "terraform"
	]
	privileges {
	  type   = "Create"
	  target = "Policy"
	  scope {
		all = true
	  }
	}
	privileges {
	  type         = "Create"
	  target       = "Entitlement"
	  default_tags = ["cc"]
	}
	privileges {
	  type      = "AssignFunction"
	  target    = "All"
	  functions = ["Connector", "Controller", "GateWAY", "logserver"]
	}
}
`, context)
}

func testAccCheckadministrativeRoleWtihAssignFunctionUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_administrative_role" "test_administrative_role" {
	name  = "%{name}"
	notes =  "hello"
	tags = [
	  "terraform"
	]
	privileges {
	  type      = "AssignFunction"
	  target    = "All"
	  functions = ["Connector", "Controller", "GateWAY", "logserver"]
	}
}
`, context)
}
