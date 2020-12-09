package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateApplianceCustomizationDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                resource "appgate_appliance_customization" "test_appliance_customization" {
                    name = "%s"
                    file = "test-fixtures/appliance_customization_file.zip"

                    tags = [
                      "terraform",
                      "api-created"
                    ]
                }
                data "appgate_appliance_customization" "test" {
                    depends_on = [
                        appgate_appliance_customization.test_appliance_customization,
                    ]
                    appliance_customization_name = "%s"
                }
                `, rName, rName),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_appliance_customization.test", "appliance_customization_name"),
					resource.TestCheckResourceAttrSet("data.appgate_appliance_customization.test", "appliance_customization_id"),
				),
			},
		},
	})
}
