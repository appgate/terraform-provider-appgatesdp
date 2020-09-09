package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDeviceScriptBasic(t *testing.T) {
	resourceName := "appgate_device_script.test_device_script"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeviceScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDeviceScriptBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceScriptExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "checksum_sha256", "74443048b52bf2be3b0f003a8f37592551d316c6c7c28b1a110ee6f879ef4130"),
					resource.TestCheckResourceAttr(resourceName, "content", "#!/usr/bin/env bash\necho \"hello world\"\n"),

					resource.TestCheckResourceAttr(resourceName, "filename", "acceptance_script.sh"),
					resource.TestCheckResourceAttr(resourceName, "name", "device_script_one"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2876187004", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccDeviceScriptImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckDeviceScriptBasic() string {
	return fmt.Sprintf(`
resource "appgate_device_script" "test_device_script" {
  name     = "device_script_one"
  filename = "acceptance_script.sh"
  content  = <<-EOF
#!/usr/bin/env bash
echo "hello world"
EOF
  tags = [
    "terraform",
    "api-created"
  ]
}
`)
}

func testAccCheckDeviceScriptExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.DeviceScriptsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.DeviceScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching device script with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckDeviceScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_device_script" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.DeviceScriptsApi

		_, _, err := api.DeviceScriptsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Device script still exists, %+v", err)
		}
	}
	return nil
}

func testAccDeviceScriptImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
