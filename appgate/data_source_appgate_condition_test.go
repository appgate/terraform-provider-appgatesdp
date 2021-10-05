package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateConditionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: conditionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgatesdp_condition.test_data_condition", "condition_id"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_condition.test_data_condition", "condition_name"),
					resource.TestCheckResourceAttr("data.appgatesdp_condition.test_data_condition", "condition_name", "Always"),
				),
			},
		},
	})
}

func conditionDataSourceConfig() string {
	return `
data "appgatesdp_condition" "test_data_condition" {
    condition_name = "Always"
}`
}
