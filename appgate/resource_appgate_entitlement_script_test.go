package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEntitlementScriptBasic(t *testing.T) {
	resourceName := "appgatesdp_entitlement_script.test_entitlement_script"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEntitlementScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementScriptBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", "return [];"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "appShortcut"),
					resource.TestCheckResourceAttr(resourceName, "notes", "test only"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccEntitlementScriptImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckEntitlementScriptBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_entitlement_script" "test_entitlement_script" {
  name       = "%s"
  type       = "appShortcut"
  expression = "return [];"
  notes      = "test only"
  tags = [
    "terraform",
    "api-created"
  ]
}
`, rName)
}

func testAccCheckEntitlementScriptExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}

		api := testAccProvider.Meta().(*Client).API.EntitlementScriptsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.EntitlementScriptsIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err != nil {
			return fmt.Errorf("error fetching entitlement script with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckEntitlementScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_entitlement_script" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.EntitlementScriptsApi

		if _, _, err := api.EntitlementScriptsIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("Criteria script still exists, %+v", err)
		}
	}
	return nil
}

func testAccEntitlementScriptImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
