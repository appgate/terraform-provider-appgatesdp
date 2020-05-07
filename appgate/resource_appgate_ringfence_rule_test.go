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
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.direction", "out"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.hosts.0", "10.0.2.0/24"),

					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.ports.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.ports.0", "80"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.ports.1", "443"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.ports.2", "1024-2048"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4182075683.types.0", "0-255"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2876187004", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
				),
			},
		},
	})
}

func testAccCheckRingfenceRule() string {
	return fmt.Sprintf(`
resource "appgate_ringfence_rule" "test_ringfence_rule" {
    name = "ringfence-rule-test"
    tags = [
        "terraform",
        "api-created"
      ]

      actions {
        protocol  = "icmp"
        direction = "out"
        action    = "allow"

        hosts = [
          "10.0.2.0/24"
        ]

        ports = [
          "80",
          "443",
          "1024-2048"
        ]

        types = [
          "0-255"
        ]

      }
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
