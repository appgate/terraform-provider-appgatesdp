package appgate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppgateEntitlementDataSource(t *testing.T) {
	dataSourceNameAB := "data.appgatesdp_entitlement.test_get_ab"
	resourceNameAB := "appgatesdp_entitlement.test_ab"
	dataSourceNameB := "data.appgatesdp_entitlement.test_get_b"
	resourceNameB := "appgatesdp_entitlement.test_b"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testEntitlementDataSourcCreateResources(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntitlementExists(resourceNameAB),
					testAccCheckEntitlementExists(resourceNameB),
					resource.TestCheckResourceAttr(resourceNameAB, "name", "ab"),
					resource.TestCheckResourceAttr(resourceNameB, "name", "b"),
				),
			},
			{
				ResourceName:      resourceNameAB,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccEntitlementImportStateCheckFunc(1),
			},
			{
				Config: testEntitlementDataSourceSimilarNamePrefixes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAB, "name", "ab"),
					resource.TestCheckResourceAttrPair(dataSourceNameAB, "entitlement_name", resourceNameAB, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAB, "entitlement_id", resourceNameAB, "id"),
				),
			},
			{
				Config: testEntitlementDataSourceSimilarFirstEnttitlment(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameB, "name", "b"),
					resource.TestCheckResourceAttrPair(dataSourceNameB, "entitlement_name", resourceNameB, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameB, "entitlement_id", resourceNameB, "id"),
				),
			},
			{
				Config: testEntitlementDataSourceFromID(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAB, "name", "ab"),
					resource.TestCheckResourceAttrPair(dataSourceNameAB, "entitlement_name", resourceNameAB, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameAB, "entitlement_id", resourceNameAB, "id"),
				),
			},
		},
	})
}

func testEntitlementDataSourcCreateResources() string {
	return `
data "appgatesdp_site" "default_site" {
	site_name = "Default Site"
}
data "appgatesdp_condition" "always" {
	condition_name = "Always"
}
resource "appgatesdp_entitlement" "test_aba" {
	name = "aba"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created",
		"b1"
	]
	disabled = true
	condition_logic = "and"
	actions {
		subtype = "icmp_up"
		action  = "allow"
		types   = ["0-16"]
		hosts = [
			"10.0.0.1",
		]
	}
}	
resource "appgatesdp_entitlement" "test_b" {
	name = "b"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created",
		"b1"
	]
	disabled = true
	condition_logic = "and"
	actions {
		subtype = "icmp_up"
		action  = "allow"
		types   = ["0-16"]
		hosts = [
			"10.0.0.1",
		]
	}
}
resource "appgatesdp_entitlement" "test_ab" {
	name = "ab"
	site = data.appgatesdp_site.default_site.id
	conditions = [
		data.appgatesdp_condition.always.id
	]
	tags = [
		"terraform",
		"api-created",
		"ab1"
	]
	disabled = true
	condition_logic = "and"
	actions {
		subtype = "icmp_up"
		action  = "allow"
		types   = ["0-16"]
		hosts = [
			"10.0.0.1",
		]
	}
}
	`
}

func testEntitlementDataSourceSimilarNamePrefixes() string {
	return testEntitlementDataSourcCreateResources() + `
data "appgatesdp_entitlement" "test_get_ab" {
	entitlement_name = "ab"
}
`
}

func testEntitlementDataSourceSimilarFirstEnttitlment() string {
	return testEntitlementDataSourcCreateResources() + `
data "appgatesdp_entitlement" "test_get_b" {
	entitlement_name = "b"
}
`
}

func testEntitlementDataSourceFromID() string {
	return testEntitlementDataSourcCreateResources() + `
data "appgatesdp_entitlement" "test_get_ab" {
	entitlement_id = appgatesdp_entitlement.test_ab.id
}
`
}
