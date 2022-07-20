package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateUserClaimScriptataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgatesdp_user_claim_script.test"
	resourceName := "appgatesdp_user_claim_script.test_user_claim_script"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
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
                data "appgatesdp_user_claim_script" "test" {
                    user_claim_script_id = appgatesdp_user_claim_script.test_user_claim_script.id
                }
                `, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "user_claim_script_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_claim_script_id", resourceName, "id"),
				),
			},
		},
	})
}
