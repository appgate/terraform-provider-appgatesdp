package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateIPPoolDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                resource "appgate_ip_pool" "test_data_ip_pool" {
                    name            = "%s"
                    lease_time_days = 5
                    ranges {
                      first = "10.0.0.1"
                      last  = "10.0.0.254"
                    }

                    tags = [
                      "terraform",
                      "api-created"
                    ]
                }
				data "appgate_ip_pool" "test_ip_pool_data_source" {
                    depends_on = [
                        appgate_ip_pool.test_data_ip_pool,
                    ]
                    ip_pool_name = "%s"
                }
                `, rName, rName),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_ip_pool.test_ip_pool_data_source", "ip_pool_name"),
					resource.TestCheckResourceAttrSet("data.appgate_ip_pool.test_ip_pool_data_source", "ip_pool_id"),
				),
			},
		},
	})
}
