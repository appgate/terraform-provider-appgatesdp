package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCriteriaScriptBasic(t *testing.T) {
	resourceName := "appgatesdp_criteria_script.test_criteria_script"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCriteriaScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCriteriaScriptBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCriteriaScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "expression", "return claims.user.username === 'admin';"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "bb"),
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

func testAccCheckCriteriaScriptBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_criteria_script" "test_criteria_script" {
  name       = "%s"
  expression = "return claims.user.username === 'admin';"
  tags = [
    "aa",
    "bb"
  ]
}
`, rName)
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
		if rs.Type != "appgatesdp_criteria_script" {
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
