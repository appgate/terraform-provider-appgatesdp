package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateApplianceCustomizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                resource "appgate_appliance_customization" "test_appliance_customization" {
                    name = "test customization"
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
                    appliance_customization_name = "test customization"
                }
                `,
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
