package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyStopBasic(t *testing.T) {
	resourceName := "appgatesdp_stop_policy.test_stop_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":     rName,
		"new_name": rName + "NEW",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testFor62AndAbove(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyStopBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "aa"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "bb"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "cc"),
					resource.TestCheckResourceAttr(resourceName, "type", "Stop"),
					resource.TestCheckResourceAttr(resourceName, "client_profile_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_profile_settings.0.force", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_profile_settings.0.profiles.#", "0"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyStopBasicUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expression", emptyPolicyExpression),
					resource.TestCheckResourceAttr(resourceName, "name", context["new_name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "hello"),
					resource.TestCheckResourceAttr(resourceName, "type", "Stop"),
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

func testAccCheckPolicyStopBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_stop_policy" "test_stop_policy" {
	name = "%{name}"
	tags = ["aa", "bb", "cc"]
	client_profile_settings {
		enabled = true
        force = true
		profiles = []
	}
}
`, context)
}

func testAccCheckPolicyStopBasicUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_stop_policy" "test_stop_policy" {
	name = "%{new_name}"
	tags = ["hello"]
	client_profile_settings {
		enabled = true
        force = true
		profiles = []
	}
}
`, context)
}
