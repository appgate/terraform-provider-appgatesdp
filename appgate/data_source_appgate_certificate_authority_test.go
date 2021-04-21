package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateCADataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `data "appgatesdp_certificate_authority" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "version"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "serial"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "issuer"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "subject"),

					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "valid_from"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "valid_to"),

					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "fingerprint"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "certificate"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.test", "subject_public_key"),
				),
			},
		},
	})
}
