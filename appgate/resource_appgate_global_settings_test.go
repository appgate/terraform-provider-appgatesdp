package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGlobalSettingsBasic(t *testing.T) {
	resourceName := "appgatesdp_global_settings.test_global_settings"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGlobalSettingsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "administration_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "app_discovery_domains.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_log_persistence_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_api_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "collective_id"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_token_expiration", "500"),
					resource.TestCheckResourceAttrSet(resourceName, "fips"),
					resource.TestCheckResourceAttrSet(resourceName, "geo_ip_updates"),
					resource.TestCheckResourceAttr(resourceName, "login_banner_message", "Welcome"),
					resource.TestCheckResourceAttr(resourceName, "message_of_the_day", "hello world"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_certificate_expiration"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccGlobalSettingsImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckGlobalSettingsBasic() string {
	return `
resource "appgatesdp_global_settings" "test_global_settings" {
  message_of_the_day           = "hello world"
  entitlement_token_expiration = 500
  login_banner_message         = "Welcome"
}
`
}

func testAccCheckGlobalSettingsExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.GlobalSettingsApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.GlobalSettingsGet(context.Background()).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching global settings with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccGlobalSettingsImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}

func TestAccGlobalSettings54ProfileHostname(t *testing.T) {
	resourceName := "appgatesdp_global_settings.test_global_settings_profile_hostname"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.LessThan(Appliance54Version) {
						t.Skipf("Test only for 5.4 and above, client_connections profile_hostname is not supported prior to 5.4, you are using %s", currentVersion.String())
					}
				},
				Config: testAccGlobalSettingsProfileHostname(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "administration_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "app_discovery_domains.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_log_persistence_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_api_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "collective_id"),
					resource.TestCheckResourceAttrSet(resourceName, "entitlement_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "fips"),
					resource.TestCheckResourceAttrSet(resourceName, "geo_ip_updates"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_certificate_expiration"),
					resource.TestCheckResourceAttr(resourceName, "profile_hostname", "xyz.appgate-sdp.com"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccGlobalSettingsImportStateCheckFunc(1),
			},
		},
	})
}

func testAccGlobalSettingsProfileHostname() string {
	return `
resource "appgatesdp_global_settings" "test_global_settings_profile_hostname" {
	profile_hostname = "xyz.appgate-sdp.com"
}
`
}

func TestAccGlobalSettings6AndAbove(t *testing.T) {
	resourceName := "appgatesdp_global_settings.test_global_settings"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor6AndAbove(t)
				},
				Config: `
				resource "appgatesdp_global_settings" "test_global_settings" {
					message_of_the_day           = "hello world"
					entitlement_token_expiration = 500
					login_banner_message         = "Welcome"
					collective_name              = "Tardis"
					spa_time_window_seconds      = 200
					spa_mode                     = "TCP"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "administration_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "app_discovery_domains.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_log_persistence_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_api_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "collective_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fips"),
					resource.TestCheckResourceAttrSet(resourceName, "geo_ip_updates"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_certificate_expiration"),

					resource.TestCheckResourceAttr(resourceName, "entitlement_token_expiration", "500"),
					resource.TestCheckResourceAttr(resourceName, "login_banner_message", "Welcome"),
					resource.TestCheckResourceAttr(resourceName, "message_of_the_day", "hello world"),
					resource.TestCheckResourceAttr(resourceName, "collective_name", "Tardis"),
					resource.TestCheckResourceAttr(resourceName, "spa_time_window_seconds", "200"),
					resource.TestCheckResourceAttr(resourceName, "spa_mode", "TCP"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccGlobalSettingsImportStateCheckFunc(1),
			},
			{
				PreConfig: func() {
					testFor6AndAbove(t)
				},
				Config: `
				resource "appgatesdp_global_settings" "test_global_settings" {
					message_of_the_day           = "Welcome to appgate sdp"
					entitlement_token_expiration = 400
					login_banner_message         = "Hello world"
					collective_name              = "Who"
					spa_time_window_seconds      = 233
					spa_mode                     = "UDP-TCP"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(resourceName),
					// testAccCheckExampleWidgetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "administration_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "app_discovery_domains.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_log_persistence_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_api_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "collective_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fips"),
					resource.TestCheckResourceAttrSet(resourceName, "geo_ip_updates"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_certificate_expiration"),

					resource.TestCheckResourceAttr(resourceName, "entitlement_token_expiration", "400"),
					resource.TestCheckResourceAttr(resourceName, "login_banner_message", "Hello world"),
					resource.TestCheckResourceAttr(resourceName, "message_of_the_day", "Welcome to appgate sdp"),
					resource.TestCheckResourceAttr(resourceName, "collective_name", "Who"),
					resource.TestCheckResourceAttr(resourceName, "spa_time_window_seconds", "233"),
					resource.TestCheckResourceAttr(resourceName, "spa_mode", "UDP-TCP"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccGlobalSettingsImportStateCheckFunc(1),
			},
		},
	})
}

func TestAccGlobalSettings62AndAbove(t *testing.T) {
	resourceName := "appgatesdp_global_settings.test_global_settings"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor62AndAbove(t)
				},
				Config: `
				resource "appgatesdp_global_settings" "test_global_settings" {
					message_of_the_day           = "hello world"
					entitlement_token_expiration = 500
					login_banner_message         = "Welcome"
					collective_name              = "Tardis"
					spa_time_window_seconds      = 222
					registered_device_expiration_days = 111
					spa_mode                     = "TCP"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "administration_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "app_discovery_domains.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_log_persistence_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "backup_api_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_token_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "collective_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fips"),
					resource.TestCheckResourceAttrSet(resourceName, "geo_ip_updates"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_certificate_expiration"),

					resource.TestCheckResourceAttr(resourceName, "entitlement_token_expiration", "500"),
					resource.TestCheckResourceAttr(resourceName, "login_banner_message", "Welcome"),
					resource.TestCheckResourceAttr(resourceName, "message_of_the_day", "hello world"),
					resource.TestCheckResourceAttr(resourceName, "collective_name", "Tardis"),
					resource.TestCheckResourceAttr(resourceName, "spa_time_window_seconds", "222"),
					resource.TestCheckResourceAttr(resourceName, "registered_device_expiration_days", "111"),
					resource.TestCheckResourceAttr(resourceName, "spa_mode", "TCP"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccGlobalSettingsImportStateCheckFunc(1),
			},
		},
	})
}
