package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateGlobalSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `data "appgatesdp_global_settings" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "claims_token_expiration"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "entitlement_token_expiration"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "administration_token_expiration"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "vpn_certificate_expiration"),

					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "backup_api_enabled"),

					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "fips"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "geo_ip_updates"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "audit_log_persistence_mode"),

					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.test", "collective_id"),
				),
			},
		},
	})
}
