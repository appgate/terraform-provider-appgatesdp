package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateApplianceCustomizationDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgate_appliance_customization.test"
	resourceName := "appgate_appliance_customization.test_appliance_customization"
	resource.ParallelTest(t, resource.TestCase{
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
                    appliance_customization_id = appgate_appliance_customization.test_appliance_customization.id
                }
                `, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "appliance_customization_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "appliance_customization_id", resourceName, "id"),
				),
			},
		},
	})
}
