package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateIPPoolDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgate_ip_pool.test_ip_pool_data_source"
	resourceName := "appgate_ip_pool.test_data_ip_pool"
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
                    ip_pool_id = appgate_ip_pool.test_data_ip_pool.id
                }
                `, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_pool_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_pool_id", resourceName, "id"),
				),
			},
		},
	})
}
