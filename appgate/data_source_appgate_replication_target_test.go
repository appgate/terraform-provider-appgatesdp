package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateReplicationTargetDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgatesdp_replication_target.testdd"
	resourceName := "appgatesdp_replication_target.test_data_replication_target"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "appgatesdp_replication_target" "test_data_replication_target" {
				  name = "%s"
				  replication_tags = ["test"]
				}
				data "appgatesdp_replication_target" "testdd" {
                    replication_target_id = appgatesdp_replication_target.test_data_replication_target.id
                }
                `, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "replication_target_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "replication_target_id", resourceName, "id"),
				),
			},
		},
	})
}
