package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRingfenceRuleBasicICMP(t *testing.T) {
	resourceName := "appgate_ringfence_rule.test_ringfence_rule"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRingfenceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRingfenceRuleICMP(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRingfenceRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.direction", "out"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "10.0.2.0/24"),

					resource.TestCheckResourceAttr(resourceName, "actions.0.protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
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

func testAccCheckRingfenceRuleICMP(rName string) string {
	return fmt.Sprintf(`
resource "appgate_ringfence_rule" "test_ringfence_rule" {
    name = "%s"
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
        types = [
          "0-255"
        ]
      }
}
`, rName)
}

func TestAccRingfenceRuleBasicTCP(t *testing.T) {
	resourceName := "appgate_ringfence_rule.test_ringfence_rule_tcp"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRingfenceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRingfenceRuleTCP(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRingfenceRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.direction", "out"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.hosts.0", "10.0.2.0/24"),

					resource.TestCheckResourceAttr(resourceName, "actions.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.0", "80"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.1", "443"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.ports.2", "1024-2048"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
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
func testAccCheckRingfenceRuleTCP(rName string) string {
	return fmt.Sprintf(`
resource "appgate_ringfence_rule" "test_ringfence_rule_tcp" {
    name = "%s"
    tags = [
      "terraform",
      "api-created"
    ]
    actions {
      protocol = "tcp"
      action   = "allow"
      direction = "out"

      hosts = [
        "10.0.2.0/24"
      ]

      ports = [
        "80",
        "443",
        "1024-2048"
      ]
    }
}
`, rName)
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
