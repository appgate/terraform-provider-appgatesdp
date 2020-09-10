package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateIPPoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                resource "appgate_ip_pool" "test_ip_pool" {
                    name            = "ip pool test"
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
				data "appgate_ip_pool" "test" {
                    depends_on = [
                        appgate_ip_pool.test_ip_pool,
                    ]
                    ip_pool_name = "ip pool test"
                }
                `,
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_ip_pool.test", "ip_pool_name"),
					resource.TestCheckResourceAttrSet("data.appgate_ip_pool.test", "ip_pool_id"),
				),
			},
		},
	})
}
