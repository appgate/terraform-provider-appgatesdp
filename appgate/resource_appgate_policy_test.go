package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPolicyBasic(t *testing.T) {
	resourceName := "appgatesdp_policy.test_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "expression", "return true;\n"),
					resource.TestCheckResourceAttr(resourceName, "notes", "terraform policy notes"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.0", "developer"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.0", "developer"),
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

func testAccCheckPolicyBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_policy" "test_policy" {
    name  = "%s"
    notes = "terraform policy notes"
    tags = [
        "terraform",
        "api-created"
    ]
    disabled = false

    expression = <<-EOF
        return true;
    EOF
    entitlement_links = [
        "developer"
    ]
    ringfence_rule_links = [
        "developer"
    ]
    tamper_proofing = true
    proxy_auto_config {
        enabled = true
        url     = "http://foo.com"
        persist = false
    }
    trusted_network_check {
        enabled    = true
        dns_suffix = "aa"
    }
}

`, rName)
}

func testAccCheckPolicyExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.PoliciesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.PoliciesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching policy with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_policy" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.PoliciesApi

		if _, _, err := api.PoliciesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("policy still exists, %+v", err)
		}
	}
	return nil
}
