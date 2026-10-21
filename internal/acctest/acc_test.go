package acctest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/appgate/terraform-provider-appgatesdp/internal/acctest"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/provider"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

var factories = map[string]func() (tfprotov6.ProviderServer, error){
	"appgatesdp": providerserver.NewProtocol6WithError(provider.New("acc")()),
}

// TestAccSchemaConformance is read-only: it lists objects of every
// supported collection and singleton and compares them with the generated
// metadata, surfacing differences between the specification and the API.
func TestAccSchemaConformance(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		c, api := env.Client(t)
		ctx := context.Background()
		for _, name := range sortedNames(api.DataSources) {
			r := api.DataSources[name]
			t.Run(name, func(t *testing.T) {
				var objects []map[string]any
				if r.Singleton {
					var obj map[string]any
					if err := c.Do(ctx, http.MethodGet, r.Path, nil, nil, &obj); err != nil {
						t.Fatal(err)
					}
					objects = append(objects, obj)
				} else {
					var list struct {
						Data []map[string]any `json:"data"`
					}
					if err := c.Do(ctx, http.MethodGet, r.Path, url.Values{"range": {"0-20"}}, nil, &list); err != nil {
						t.Fatal(err)
					}
					objects = list.Data
				}
				for _, obj := range objects {
					if r.Discriminator != nil && obj[r.Discriminator.Property] != r.Discriminator.Value {
						continue
					}
					for _, f := range acctest.Conformance(fmt.Sprintf("%s[%v]", name, obj["id"]), obj, r.Fields) {
						if f.Severity == "error" {
							t.Error(f)
						} else {
							t.Log(f)
						}
					}
				}
			})
		}
	})
}

// TestAccConditionLifecycle creates, updates, imports and destroys a
// condition in every configured environment.
func TestAccConditionLifecycle(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		name := acctest.Name("condition")
		cfg := func(expr string) string {
			return env.ProviderConfig() + fmt.Sprintf(`
resource "appgatesdp_condition" "test" {
  name       = %q
  expression = %q
  tags       = ["terraform-acceptance"]
}
`, name, expr)
		}
		sweepAfter(t, env)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			Steps: []resource.TestStep{
				{Config: cfg("return true;"), Check: resource.TestCheckResourceAttr("appgatesdp_condition.test", "name", name)},
				{Config: cfg("return false;"), Check: resource.TestCheckResourceAttr("appgatesdp_condition.test", "expression", "return false;")},
				{ResourceName: "appgatesdp_condition.test", ImportState: true, ImportStateVerify: true},
			},
		})
	})
}

// TestAccPolicyAndEntitlement covers references between resources and a
// discriminated collection.
func TestAccPolicyAndEntitlement(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		sweepAfter(t, env)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			Steps: []resource.TestStep{{
				Config: env.ProviderConfig() + fmt.Sprintf(`
resource "appgatesdp_site" "test" {
  name = %q
}
resource "appgatesdp_entitlement" "test" {
  name       = %q
  site       = appgatesdp_site.test.id
  conditions = []
  actions = [{
    subtype = "tcp_up"
    action  = "allow"
    hosts   = ["10.0.0.1"]
    ports   = ["443"]
  }]
}
resource "appgatesdp_access_policy" "test" {
  name         = %q
  expression   = "return true;"
  entitlements = [appgatesdp_entitlement.test.id]
}
`, acctest.Name("site"), acctest.Name("entitlement"), acctest.Name("policy")),
				Check: resource.TestCheckResourceAttrPair("appgatesdp_access_policy.test", "entitlements.0", "appgatesdp_entitlement.test", "id"),
			}},
		})
	})
}

// TestAccBlacklistUser blacklists a fictitious user and removes the entry.
func TestAccBlacklistUser(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		sweepAfter(t, env)
		dn := fmt.Sprintf("CN=%s,OU=local", acctest.Name("user"))
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			Steps: []resource.TestStep{
				{
					Config: env.ProviderConfig() + fmt.Sprintf(`
resource "appgatesdp_blacklist_user" "test" {
  user_distinguished_name = %q
  reason                  = "terraform acceptance test"
}
`, dn),
					Check: resource.TestCheckResourceAttr("appgatesdp_blacklist_user.test", "id", dn),
				},
				{ResourceName: "appgatesdp_blacklist_user.test", ImportState: true, ImportStateId: dn, ImportStateVerify: true},
			},
		})
	})
}

// TestAccSweep removes leftovers of interrupted acceptance runs.
func TestAccSweep(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		c, api := env.Client(t)
		if err := acctest.Sweep(context.Background(), c, api); err != nil {
			t.Fatal(err)
		}
	})
}

func sweepAfter(t *testing.T, env acctest.Env) {
	c, api := env.Client(t)
	t.Cleanup(func() {
		if err := acctest.Sweep(context.Background(), c, api); err != nil {
			t.Errorf("cleanup: %v", err)
		}
	})
}

func sortedNames[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

// TestAcceptanceSuiteAgainstMock runs the acceptance suite against the mock
// controller, so the suite itself is verified without a live environment.
func TestAcceptanceSuiteAgainstMock(t *testing.T) {
	api, err := generated.API("v25")
	if err != nil {
		t.Fatal(err)
	}
	srv := sdptest.New(t, api)
	srv.Seed("/conditions", map[string]any{"name": "existing", "expression": "return true;"})
	t.Setenv("TF_ACC", "1")
	t.Setenv("SDP_ACC_ALLOW_DESTRUCTIVE", "1")
	t.Setenv("SDP_ACC_V25_URL", srv.URL)
	t.Setenv("SDP_ACC_V25_USERNAME", sdptest.Username)
	t.Setenv("SDP_ACC_V25_PASSWORD", sdptest.Password)
	for name, test := range map[string]func(*testing.T){
		"conformance": TestAccSchemaConformance,
		"condition":   TestAccConditionLifecycle,
		"entitlement": TestAccPolicyAndEntitlement,
		"sweep":       TestAccSweep,
	} {
		t.Run(name, test)
	}
	for _, path := range []string{"/conditions", "/sites", "/entitlements", "/policies", "/blacklist"} {
		if n := srv.Count(path); path != "/conditions" && n != 0 || path == "/conditions" && n != 1 {
			t.Errorf("%s holds %d objects after the suite; test objects must be cleaned up", path, n)
		}
	}
}
