package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccEntitlementScriptBasic(t *testing.T) {
	resourceName := "appgate_entitlement_script.test_entitlement_script"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEntitlementScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEntitlementScriptBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", "return [];"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_entitlement_script"),
					resource.TestCheckResourceAttr(resourceName, "notes", "test only"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2876187004", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
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

func testAccCheckEntitlementScriptBasic() string {
	return fmt.Sprintf(`
resource "appgate_entitlement_script" "test_entitlement_script" {
  name       = "test_entitlement_script"
  expression = "return [];"
  notes      = "test only"
  tags = [
    "terraform",
    "api-created"
  ]
}
`)
}

func testAccCheckEntitlementScriptExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.EntitlementScriptsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.EntitlementScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching entitlement script with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckEntitlementScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_entitlement_script" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.EntitlementScriptsApi

		_, _, err := api.EntitlementScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
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
