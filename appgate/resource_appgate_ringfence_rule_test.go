package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRingfenceRuleBasic(t *testing.T) {
	resourceName := "appgate_ringfence_rule.test_ringfence_rule"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRingfenceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRingfenceRule(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRingfenceRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "ringfence-rule-test"),
				),
			},
		},
	})
}

func testAccCheckRingfenceRule() string {
	return fmt.Sprintf(`

resource "appgate_ringfence_rule" "test_ringfence_rule" {
    name = "ringfence-rule-test"
}
`)
}

func testAccCheckRingfenceRuleExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.RingfenceRulesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.RingfenceRulesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching ringfence rule with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckRingfenceRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_ringfence_rule" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.RingfenceRulesApi

		_, _, err := api.RingfenceRulesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("RingfenceRule still exists, %+v", err)
		}
	}
	return nil
}
