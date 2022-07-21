package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserClaimScriptBasic(t *testing.T) {
	resourceName := "appgatesdp_user_claim_script.test_user_claim_script"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserClaimScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserClaimScriptBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserClaimScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "expression", "return {'posture': 25};\n"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccUserClaimScriptImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckUserClaimScriptUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserClaimScriptExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "updated claim name"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "expression", "return {'foo': 'bar'};"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "hello"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccUserClaimScriptImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckUserClaimScriptBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_user_claim_script" "test_user_claim_script" {
  name     = "%s"
  expression  = <<-EOF
return {'posture': 25};
EOF
  tags = [
    "terraform",
    "api-created"
  ]
}
`, rName)
}

func testAccCheckUserClaimScriptUpdated() string {
	return `
resource "appgatesdp_user_claim_script" "test_user_claim_script" {
  name     = "updated claim name"
  expression = "return {'foo': 'bar'};"
  tags = [
    "hello"
  ]
}
`
}

func testAccCheckUserClaimScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_user_claim_script" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.UserClaimScriptsApi

		if _, _, err := api.UserScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("user claim script still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckUserClaimScriptExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.UserClaimScriptsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.UserScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching user claim script with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccUserClaimScriptImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
