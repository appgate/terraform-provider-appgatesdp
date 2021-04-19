package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConditionBasic(t *testing.T) {
	resourceName := "appgate_condition.test_condition"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCondition(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConditionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "expression", "return true;"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "remedy_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "remedy_methods.0.message", "This resource requires you to enter your password again"),
					resource.TestCheckResourceAttr(resourceName, "remedy_methods.0.type", "DisplayMessage"),

					resource.TestCheckResourceAttr(resourceName, "repeat_schedules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repeat_schedules.0", "13:32"),
					resource.TestCheckResourceAttr(resourceName, "repeat_schedules.1", "1h"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckCondition(rName string) string {
	return fmt.Sprintf(`
resource "appgate_condition" "test_condition" {
    name = "%s"
    tags = [
      "terraform",
      "api-created"
    ]

    expression = "return true;"

    repeat_schedules = [
      "1h",
      "13:32"
    ]
    remedy_methods {
        type        = "DisplayMessage"
        message     = "This resource requires you to enter your password again"
    }
}
`, rName)
}

func testAccCheckConditionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ConditionsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.ConditionsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching condition with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckConditionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_condition" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ConditionsApi

		_, _, err := api.ConditionsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Condition still exists, %+v", err)
		}
	}
	return nil
}
