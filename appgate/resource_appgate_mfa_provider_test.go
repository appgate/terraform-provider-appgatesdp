package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMfaProviderBasic(t *testing.T) {
	resourceName := "appgatesdp_mfa_provider.test_mfa_provider"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMfaProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMfaProviderBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMfaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authentication_protocol", "CHAP"),
					resource.TestCheckResourceAttr(resourceName, "challenge_shared_secret", "secretString"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "mfa.company.com"),
					resource.TestCheckResourceAttr(resourceName, "mode", "Challenge"),
					resource.TestCheckResourceAttr(resourceName, "name", "themfaprovider"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "port", "1812"),
					resource.TestCheckResourceAttr(resourceName, "shared_secret", "helloworld"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "6"),
					resource.TestCheckResourceAttr(resourceName, "type", "Radius"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccMfaProviderImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckMfaProviderBasic() string {
	return `
resource "appgatesdp_mfa_provider" "test_mfa_provider" {
  name = "themfaprovider"
  port = 1812
  type = "Radius"
	shared_secret = "helloworld"
	challenge_shared_secret = "secretString"
  hostnames = [
    "mfa.company.com"
  ]

  tags = [
    "terraform",
    "api-created"
  ]
}
`
}

func testAccCheckMfaProviderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.MFAProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.MfaProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching mfa_provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckMfaProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_mfa_provider" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.MFAProvidersApi

		_, _, err := api.MfaProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("mfa_provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccMfaProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
