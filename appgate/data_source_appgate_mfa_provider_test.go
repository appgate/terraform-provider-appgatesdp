package appgate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateMfaProviderDataSource(t *testing.T) {
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	dataSourceName := "data.appgate_mfa_provider.test_mfa_ds"
	resourceName := "appgate_mfa_provider.test_mfa_provider"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "appgate_mfa_provider" "test_mfa_provider" {
					  name                    = "%s"
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
					data "appgate_mfa_provider" "test_mfa_ds" {
					  mfa_provider_id = appgate_mfa_provider.test_mfa_provider.id
					}
                `, rName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "mfa_provider_name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mfa_provider_id", resourceName, "id"),
				),
			},
		},
	})
}
