package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
)

// TestConditionLifecycle runs create, update, import, drift correction,
// out-of-band deletion and destroy against every API version.
func TestConditionLifecycle(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		srv := newServer(t, v.APIVersion)
		srv.LowercaseTags = true
		cfg := func(expr string) string {
			return providerConfig(srv, v.APIVersion) + fmt.Sprintf(`
resource "appgatesdp_condition" "test" {
  name       = "tf-condition"
  expression = %q
  tags       = ["Prod", "api"]
}
`, expr)
		}
		var firstID string
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: protoV6Factories,
			CheckDestroy: func(*terraform.State) error {
				if n := srv.Count("/conditions"); n != 0 {
					return fmt.Errorf("%d conditions left after destroy", n)
				}
				return nil
			},
			Steps: []resource.TestStep{
				{
					Config: cfg("return true;"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("appgatesdp_condition.test", "id"),
						resource.TestCheckResourceAttr("appgatesdp_condition.test", "expression", "return true;"),
						// The controller lowercases and sorts tags; configured values are kept.
						resource.TestCheckResourceAttr("appgatesdp_condition.test", "tags.0", "Prod"),
						resource.TestCheckResourceAttr("appgatesdp_condition.test", "created", "2026-01-01T00:00:00Z"),
						checkAcceptHeaders(srv, v),
						func(s *terraform.State) (err error) {
							firstID, err = resourceID(s, "appgatesdp_condition.test")
							return err
						},
					),
				},
				{
					Config: cfg("return false;"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("appgatesdp_condition.test", "expression", "return false;"),
						func(s *terraform.State) error {
							body := lastRequest(srv, "PUT", "/conditions/")
							if body["expression"] != "return false;" || body["name"] != "tf-condition" {
								return fmt.Errorf("PUT body must contain the whole object, got %v", body)
							}
							if id, _ := resourceID(s, "appgatesdp_condition.test"); id != firstID {
								return fmt.Errorf("update replaced the object: %s != %s", id, firstID)
							}
							return nil
						},
					),
				},
				{
					ResourceName:      "appgatesdp_condition.test",
					ImportState:       true,
					ImportStateVerify: true,
					// Imported state carries the API's normalized tags.
					ImportStateVerifyIgnore: []string{"tags"},
				},
				{
					// Drift: the expression is changed outside Terraform and corrected.
					PreConfig: func() {
						srv.Update("/conditions", firstID, func(o map[string]any) { o["expression"] = "return 1;" })
					},
					Config: cfg("return false;"),
					Check: func(*terraform.State) error {
						if got := srv.Object("/conditions", firstID)["expression"]; got != "return false;" {
							return fmt.Errorf("drift not corrected, expression is %v", got)
						}
						return nil
					},
				},
				{
					// Deleted outside Terraform: read removes it and apply recreates it.
					PreConfig: func() { srv.Delete("/conditions", firstID) },
					Config:    cfg("return false;"),
					Check: func(s *terraform.State) error {
						id, err := resourceID(s, "appgatesdp_condition.test")
						if err != nil {
							return err
						}
						if id == firstID || srv.Object("/conditions", id) == nil {
							return fmt.Errorf("condition was not recreated (id %q)", id)
						}
						return nil
					},
				},
			},
		})
	})
}
