package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLocalUserBasic(t *testing.T) {
	resourceName := "appgatesdp_local_user.test_local_user"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocalUserBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", "john.doe@test.com"),
					resource.TestCheckResourceAttr(resourceName, "failed_login_attempts", "30"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "john"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "doe"),
					resource.TestCheckResourceAttr(resourceName, "lock_start", "2020-04-27T09:51:03Z"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "password", "password_is_hunter2"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+1-202-555-0172"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccLocalUserImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckLocalUserBasicUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", "jane.silly@example.com"),
					resource.TestCheckResourceAttr(resourceName, "failed_login_attempts", "30"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Jane"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Silly"),
					resource.TestCheckResourceAttr(resourceName, "lock_start", "2020-04-27T09:51:03Z"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "password", "password_is_hunter2"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+1-202-555-0172"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccLocalUserImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckLocalUserBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_local_user" "test_local_user" {
    name                  = "%s"
    first_name            = "john"
    last_name             = "doe"
    password              = "password_is_hunter2"
    email                 = "john.doe@test.com"
    phone                 = "+1-202-555-0172"
    failed_login_attempts = 30
    lock_start            = "2020-04-27T09:51:03Z"
    tags = [
      "terraform",
      "api-created"
    ]
}
`, rName)
}
func testAccCheckLocalUserBasicUpdated(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_local_user" "test_local_user" {
    name                  = "%s"
    first_name            = "Jane"
    last_name             = "Silly"
    password              = "password_is_hunter2"
    email                 = "jane.silly@example.com"
    phone                 = "+1-202-555-0172"
    failed_login_attempts = 30
    lock_start            = "2020-04-27T09:51:03Z"
    tags = [
      "terraform",
      "api-created"
    ]
}
`, rName)
}

func testAccCheckLocalUserExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LocalUsersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.LocalUsersIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err != nil {
			return fmt.Errorf("error fetching local user with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_local_user" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.LocalUsersApi

		if _, _, err := api.LocalUsersIdGet(BaseAuthContext(token), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("local user still exists, %+v", err)
		}
	}
	return nil
}

func testAccLocalUserImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
