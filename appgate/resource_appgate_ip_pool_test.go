package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIPPoolBasic(t *testing.T) {
	resourceName := "appgatesdp_ip_pool.test_ip_pool_v4"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPPoolBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_version6", "false"),
					resource.TestCheckResourceAttr(resourceName, "lease_time_days", "5"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ranges.0.first", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ranges.0.last", "10.0.0.254"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccIPPoolImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckIPPoolBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_ip_pool" "test_ip_pool_v4" {
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
`, rName)
}

func TestAccIPPoolV6(t *testing.T) {
	resourceName := "appgatesdp_ip_pool.test_ip_pool_v6"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPPoolV6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_version6", "true"),
					resource.TestCheckResourceAttr(resourceName, "lease_time_days", "5"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ranges.0.first", "2001:800:0:0:0:0:0:1"),
					resource.TestCheckResourceAttr(resourceName, "ranges.0.last", "2001:fff:ffff:ffff:ffff:ffff:ffff:fffe"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccIPPoolImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckIPPoolV6(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_ip_pool" "test_ip_pool_v6" {
	name            = "%s"
	ip_version6     = true
	lease_time_days = 5
	ranges {
	  first = "2001:800:0:0:0:0:0:1"
	  last  = "2001:fff:ffff:ffff:ffff:ffff:ffff:fffe"
	}

	tags = [
	  "terraform",
	  "api-created"
	]
}
`, rName)
}
func testAccCheckIPPoolExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.IPPoolsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.IpPoolsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching ip pool with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckIPPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_ip_pool" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.IPPoolsApi

		if _, _, err := api.IpPoolsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("Device script still exists, %+v", err)
		}
	}
	return nil
}

func testAccIPPoolImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
