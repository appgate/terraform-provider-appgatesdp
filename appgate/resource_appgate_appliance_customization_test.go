package appgate

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func TestAccApplianceCustomizationBasic(t *testing.T) {
	resourceName := "appgatesdp_appliance_customization.test_acc_appliance_customization"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	testFilename := "test-fixtures/appliance_customization_file.zip"
	testUpdatedFilename := "test-fixtures/appliance_customization_file_updated.zip"
	testFileTarget := "test-fixtures/appliance_customization_file_test.zip"

	context := map[string]interface{}{
		"name":     rName,
		"filepath": testFileTarget,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceCustomizationDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, err := copy(testFilename, testFileTarget)
					if err != nil {
						t.Errorf("Failed to copy %s to %s", testFilename, testFileTarget)
					}
				},
				Config: testAccCheckApplianceCustomizationBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceCustomizationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "checksum_sha256", "e3a9fb24832dff49ea59ff79cff9b1f24cbc0974ec62ec700165a0631fee779e"),
					resource.TestCheckResourceAttr(resourceName, "file", testFileTarget),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "size", "574"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccApplianceCustomizationImportStateCheckFunc(1),
			},
			// Test another apply with the same filename, and filepath, and make sure the sha256 checksum gets updated.
			{
				PreConfig: func() {
					_, err := copy(testUpdatedFilename, testFileTarget)
					if err != nil {
						t.Errorf("Failed to copy %s to %s", testUpdatedFilename, testUpdatedFilename)
					}
				},
				Config: testAccCheckApplianceCustomizationBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceCustomizationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "checksum_sha256", "a99ea44bff77668fea199e3d2dabe7921afb10c329453e25bf0dbc964a44606f"),
					resource.TestCheckResourceAttr(resourceName, "file", testFileTarget),
				),
			},
		},
	})
}

func testAccCheckApplianceCustomizationBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_appliance_customization" "test_acc_appliance_customization" {
    name = "%{name}"
    file = "%{filepath}"

    tags = [
      "terraform",
      "api-created"
    ]
}
`, context)
}

func testAccCheckApplianceCustomizationExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ApplianceCustomizationsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.ApplianceCustomizationsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching appliance customization with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckApplianceCustomizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_appliance_customization" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.ApplianceCustomizationsApi

		_, _, err := api.ApplianceCustomizationsIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Appliance customization still exists, %+v", err)
		}
	}
	return nil
}

func testAccApplianceCustomizationImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
