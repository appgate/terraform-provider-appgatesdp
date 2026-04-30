package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccReplication(t *testing.T) {
	resourceName := "appgatesdp_replication_target.test_replication_target"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testFor67AndAbove(t) },
		Providers:    testAccProviders,
		CheckDestroy: testReplicationCleanup,
		Steps: []resource.TestStep{
			{
				Config: testAccReplication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReplicationTargetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "replication_tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_tags.0", "repl"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccReplication(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_replication_target" "test_replication_target" {
    name             = "%s"
  	replication_tags = ["repl"]
}

resource "appgatesdp_replication_source" "test_replication_source" {
  registration_token = appgatesdp_replication_target.test_replication_target.registration_token
}
`, rName)
}

func testAccCheckReplicationTargetExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		api := testAccProvider.Meta().(*Client).API.ReplicationTargetsApi
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		if _, _, err := api.ReplicationTargetsIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err != nil {
			return fmt.Errorf("error fetching replication target with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckReplicationSourceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		api := testAccProvider.Meta().(*Client).API.ReplicationSourceApi
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		if _, _, err := api.ReplicationSourceGet(BaseAuthContext(token)).Execute(); err != nil {
			return fmt.Errorf("error fetching replication source with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testReplicationCleanup(s *terraform.State) error {
	var errors *multierror.Error
	err := testAccCheckReplicationTargetDestroy(s)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	err = testAccCheckReplicationSourceDestroy(s)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors.ErrorOrNil()
}

func testAccCheckReplicationTargetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_replication_target" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.ReplicationTargetsApi

		if _, _, err := api.ReplicationTargetsIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("ReplicationTarget still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckReplicationSourceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_replication_source" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.ReplicationSourceApi

		if _, _, err := api.ReplicationSourceGet(BaseAuthContext(token)).Execute(); err == nil {
			return fmt.Errorf("ReplicationSource still exists, %+v", err)
		}
	}
	return nil
}
