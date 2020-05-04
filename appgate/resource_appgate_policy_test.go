package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPolicyBasic(t *testing.T) {
	resourceName := "appgate_appliance.test_policy"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					// testAccCheckExampleWidgetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "policy-test"),
				),
			},
		},
	})
}

func testAccCheckPolicyBasic() string {
	return fmt.Sprintf(`
resource "appgate_policy" "test_policy" {
    name = "policy-test"
}
`)
}
