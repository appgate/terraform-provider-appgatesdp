package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateClientProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor61AndAbove(t)
				},
				Config: `
				resource "appgatesdp_client_profile" "dev_resource" {
					name                   = "dev_resource"
					spa_key_name           = "dev-resources"
					identity_provider_name = "local"
				}
				data "appgatesdp_client_profile" "dev_profile" { client_profile_id = appgatesdp_client_profile.dev_resource.id }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgatesdp_client_profile.dev_profile", "url"),
				),
			},
		},
	})
}
