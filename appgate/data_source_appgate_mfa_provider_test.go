package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAppgateMfaProviderDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "appgate_mfa_provider" "test_mfa_provider" {
					  name                    = "testMFAprovider"
					  port                    = 1812
					  type                    = "Radius"
					  shared_secret           = "helloworld"
					  challenge_shared_secret = "secretString"
					  hostnames = [
					    "mfa.company.com"
					  ]

					  tags = [
					    "terraform",
					    "api-created"
					  ]
					}
					data "appgate_mfa_provider" "test" {
					  depends_on = [
					    appgate_mfa_provider.test_mfa_provider,
					  ]
					  mfa_provider_name = "testMFAprovider"
					}
                `,
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgate_mfa_provider.test", "mfa_provider_name"),
					resource.TestCheckResourceAttrSet("data.appgate_mfa_provider.test", "mfa_provider_id"),
				),
			},
		},
	})
}
