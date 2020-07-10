package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCriteriaScriptBasic(t *testing.T) {
	resourceName := "appgate_criteria_script.test_criteria_script"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCriteriaScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCriteriaScriptBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCriteriaScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "testCriteriaScript"),
					resource.TestCheckResourceAttr(resourceName, "expression", "return claims.user.username === 'admin';"),
					resource.TestCheckResourceAttr(resourceName, "name", "testCriteriaScript"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.1389504093", "bb"),
					resource.TestCheckResourceAttr(resourceName, "tags.2075773895", "aa"),
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

func testAccCheckCriteriaScriptBasic() string {
	return fmt.Sprintf(`
resource "appgate_criteria_script" "test_criteria_script" {
  name       = "testCriteriaScript"
  expression = "return claims.user.username === 'admin';"
  tags = [
    "aa",
    "bb"
  ]
}
`)
}

func testAccCheckCriteriaScriptExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.CriteriaScriptsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.CriteriaScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching criteria script with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCriteriaScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_criteria_script" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.CriteriaScriptsApi

		_, _, err := api.CriteriaScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Criteria script still exists, %+v", err)
		}
	}
	return nil
}

func testAccCriteriaScripImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
