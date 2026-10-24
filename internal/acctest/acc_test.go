package acctest_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/appgate/terraform-provider-appgatesdp/internal/acctest"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/provider"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
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
				switch {
				case r.ReadOnly && strings.Contains(r.Path, "{"):
					t.Skip("parameterized endpoint; read by TestAccReadOnlyDataSources")
				case r.ReadOnly || r.Singleton:
					var obj map[string]any
					if err := c.DoAccept(ctx, r.Accept, http.MethodGet, r.Path, nil, nil, &obj); err != nil {
						skipUnavailable(t, err)
						t.Fatal(err)
					}
					objects = append(objects, obj)
				default:
					var list struct {
						Data []map[string]any `json:"data"`
					}
					if err := c.Do(ctx, http.MethodGet, r.Path, url.Values{"range": {"0-20"}}, nil, &list); err != nil {
						skipUnavailable(t, err)
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
  tags       = ["tf-acc-tag"]
}
`, name, expr)
		}
		sweepAfter(t, env)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			ErrorCheck:               acctest.SkipUnavailable(t),
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
			ErrorCheck:               acctest.SkipUnavailable(t),
			Steps: []resource.TestStep{{
				Config: env.ProviderConfig() + fmt.Sprintf(`
# Every collective has the default site; creating sites may exceed the
# license (sites are covered by TestAccResourceMatrix).
data "appgatesdp_site" "default" {
  name = "Default Site"
}
resource "appgatesdp_condition" "test" {
  name       = %q
  expression = "return true;"
}
resource "appgatesdp_entitlement" "test" {
  name       = %q
  site       = data.appgatesdp_site.default.id
  conditions = [appgatesdp_condition.test.id]
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
`, acctest.Name("condition"), acctest.Name("entitlement"), acctest.Name("policy")),
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
			ErrorCheck:               acctest.SkipUnavailable(t),
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

// skipUnavailable skips when the environment cannot serve a read: the
// feature is unlicensed (402), a prerequisite is missing (412), or there is
// nothing to read (404, 422, e.g. no next CA certificate). Server errors
// (5xx) and authorization errors still fail.
func skipUnavailable(t *testing.T, err error) {
	t.Helper()
	var apiErr *sdp.Error
	if errors.As(err, &apiErr) && apiErr.Status >= 400 && apiErr.Status < 500 && apiErr.Status != http.StatusUnauthorized && apiErr.Status != http.StatusForbidden {
		t.Skipf("not available in this environment: %v", err)
	}
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
	srv.Seed("/sites", map[string]any{"name": "Default Site", "networkSubnets": []any{}, "entitlementBasedRouting": true})
	// Only the mock may be targeted, whatever live environments are configured.
	for _, v := range apiversion.Default().Versions {
		t.Setenv("SDP_ACC_"+strings.ToUpper(v.APIVersion)+"_URL", "")
	}
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
		want := 0
		if path == "/conditions" || path == "/sites" {
			want = 1 // seeded objects must survive the suite
		}
		if n := srv.Count(path); n != want {
			t.Errorf("%s holds %d objects after the suite; test objects must be cleaned up", path, n)
		}
	}
}

// TestAccReplicationSourceStatus reads the replication source status
// (API v24+) and compares it with the generated metadata. Read-only.
func TestAccReplicationSourceStatus(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		c, api := env.Client(t)
		r, ok := api.Resources["replication_source"]
		if !ok {
			t.Skipf("replication source is not available in %s", env.Version.APIVersion)
		}
		var status map[string]any
		if err := c.Do(context.Background(), http.MethodGet, r.Path, nil, nil, &status); err != nil {
			skipUnavailable(t, err)
			t.Fatal(err)
		}
		for _, f := range acctest.Conformance("replication_source", status, r.Fields) {
			if f.Severity == "error" {
				t.Error(f)
			} else {
				t.Log(f)
			}
		}
	})
}

// TestAccReplicationSourceInvalidToken registers with a token the
// controller rejects. It exercises the replication source request without
// turning the collective into a replica, which a valid token would do.
func TestAccReplicationSourceInvalidToken(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		if _, api := env.Client(t); api.Resources["replication_source"] == nil {
			t.Skipf("replication source is not available in %s", env.Version.APIVersion)
		}
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			ErrorCheck:               acctest.SkipUnavailable(t),
			Steps: []resource.TestStep{{
				Config: env.ProviderConfig() + `
resource "appgatesdp_replication_source" "test" {
  token = "tf-acc-invalid-token"
}
`,
				ExpectError: regexp.MustCompile(`(?s)HTTP 422.*token: Invalid replication registration\s+token`),
			}},
		})
	})
}
