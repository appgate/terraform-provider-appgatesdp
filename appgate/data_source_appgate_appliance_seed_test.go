package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateApplianceSeedDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgatesdp_appliance_seed.test_gateway_seed_file"
	resourceName := "appgatesdp_appliance.new_test_gateway"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccSeedTest(rName),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "appliance_id", resourceName, "id"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_appliance_seed.test_gateway_seed_file", "latest_version"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_appliance_seed.test_gateway_seed_file", "password"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_appliance_seed.test_gateway_seed_file", "seed_file"),
				),
			},
		},
	})
}

func testAccSeedTest(rName string) string {
	return fmt.Sprintf(`
data "appgatesdp_site" "default_site" {
  site_name = "Default site"
}

resource "appgatesdp_appliance" "new_test_gateway" {
  name     = "%s"
  hostname = "envy-10-97-168-1337.devops"

  client_interface {
    hostname       = "envy-10-97-168-1338.devops"
    proxy_protocol = true
    https_port     = 447
    dtls_port      = 445
    allow_sources {
      address = "1.3.3.8"
      netmask = 32
      nic     = "eth0"
    }
    override_spa_mode = "UDP-TCP"
  }

  peer_interface {
    hostname   = "envy-10-97-168-1338.devops"
    https_port = "1338"

    allow_sources {
      address = "1.3.3.8"
      netmask = 32
      nic     = "eth0"
    }
  }

  site = data.appgatesdp_site.default_site.id
  networking {
    nics {
      enabled = true
      name    = "eth0"
      ipv4 {
        dhcp {
          enabled = true
          dns     = true
          routers = true
          ntp     = true
        }
      }
    }
  }

}


data "appgatesdp_appliance_seed" "test_gateway_seed_file" {
  appliance_id   = appgatesdp_appliance.new_test_gateway.id
  password       = "cz"
  latest_version = true
}
`, rName)
}
