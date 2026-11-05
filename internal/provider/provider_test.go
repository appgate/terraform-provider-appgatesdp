package provider

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

// TestProviderSchemas validates every generated schema through the
// framework's own implementation checks.
func TestProviderSchemas(t *testing.T) {
	srv, err := providerserver.NewProtocol6WithError(New("test")())()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		t.Errorf("%s: %s", d.Summary, d.Detail)
	}
	for name := range generatedResourceSchemas {
		if _, ok := resp.ResourceSchemas[TypePrefix+name]; !ok {
			t.Errorf("resource %s not served", name)
		}
	}
	for name := range generatedDataSourceSchemas {
		if _, ok := resp.DataSourceSchemas[TypePrefix+name]; !ok {
			t.Errorf("data source %s not served", name)
		}
	}
}

// TestEveryVersionHasMetadata: each registered version has generated
// metadata, and every generated resource is described by at least one version.
func TestEveryVersionHasMetadata(t *testing.T) {
	covered := map[string]bool{}
	for _, v := range apiversion.Default().Versions {
		api, err := generated.API(v.APIVersion)
		if err != nil {
			t.Fatal(err)
		}
		if api.Version != v.APIVersion {
			t.Errorf("metadata for %s claims version %s", v.APIVersion, api.Version)
		}
		for name := range api.Resources {
			if _, ok := generatedResourceSchemas[name]; !ok {
				t.Errorf("%s: metadata for %s has no schema", v.APIVersion, name)
			}
			covered[name] = true
		}
	}
	for name := range generatedResourceSchemas {
		if !covered[name] {
			t.Errorf("resource %s is not supported by any version", name)
		}
	}
	if _, err := generated.API("v1"); err == nil {
		t.Error("unknown version must fail")
	}
}

// TestEveryResourceRoundTrip creates, imports and destroys a minimal object
// of every collection resource on every version that supports it, using the
// required attributes of that version.
func TestEveryResourceRoundTrip(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		api, _ := generated.API(v.APIVersion)
		for _, name := range sortedNames(api.Resources) {
			m := api.Resources[name]
			if m.Adopt {
				continue
			}
			t.Run(name, func(t *testing.T) {
				parallelTerraform(t)
				srv := newServer(t, v.APIVersion)
				body, ok := minimalConfig(m.Fields)
				if !ok {
					t.Skip("required attributes need values the generic fixture cannot build")
				}
				addr := "appgatesdp_" + name + ".test"
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: protoV6Factories,
					Steps: []resource.TestStep{
						{
							Config: providerConfig(srv, v.APIVersion) + fmt.Sprintf("resource %q \"test\" {\n%s}\n", "appgatesdp_"+name, body),
							Check:  resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttrSet(addr, "id"), checkAcceptHeaders(srv, v)),
						},
						{ResourceName: addr, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: writeOnlyAttrs(m.Fields)},
					},
				})
			})
		}
	})
}

// minimalConfig renders the required attributes of a resource with
// placeholder values.
func minimalConfig(fields map[string]*apimeta.Field) (string, bool) {
	var out string
	for _, name := range sortedNames(fields) {
		f := fields[name]
		if !f.Required || f.ReadOnly {
			continue
		}
		var v string
		switch f.Type {
		case "string":
			v = fmt.Sprintf("%q", uuid.NewString())
			if len(f.Enum) > 0 {
				v = fmt.Sprintf("%q", f.Enum[0])
			}
		case "int64", "number":
			v = "1"
		case "bool":
			v = "true"
		case "list", "set":
			if f.Fields != nil {
				inner, ok := minimalConfig(f.Fields)
				if !ok {
					return "", false
				}
				v = "[{\n" + inner + "}]"
			} else {
				// One element satisfies minimum sizes.
				v = map[string]string{"string": fmt.Sprintf("[%q]", uuid.NewString()), "int64": "[1]", "number": "[1]", "bool": "[true]"}[f.Elem]
				if v == "" {
					return "", false
				}
			}
		case "object":
			inner, ok := minimalConfig(f.Fields)
			if !ok {
				return "", false
			}
			v = "{\n" + inner + "}"
		default:
			return "", false
		}
		out += fmt.Sprintf("  %s = %s\n", name, v)
	}
	return out, true
}

func writeOnlyAttrs(fields map[string]*apimeta.Field) []string {
	var out []string
	for name, f := range fields {
		if f.WriteOnly {
			out = append(out, name)
		}
	}
	return out
}

func TestResolveSettings(t *testing.T) {
	t.Setenv("APPGATE_ADDRESS", "https://env.example.com")
	t.Setenv("APPGATE_USERNAME", "env-user")
	t.Setenv("APPGATE_PASSWORD", "env-pass")
	t.Setenv("APPGATE_API_VERSION", "")
	t.Setenv("APPGATE_CLIENT_VERSION", "")
	t.Setenv("APPGATE_INSECURE", "")
	t.Setenv("APPGATE_DEVICE_ID", "")
	t.Setenv("APPGATE_CONFIG_PATH", "")
	t.Setenv("APPGATE_LOGIN_TIMEOUT", "")

	s, diags := resolveSettings(providerModel{Username: types.StringValue("block-user"), ClientVer: types.Int64Value(23)})
	if diags.HasError() {
		t.Fatal(diags)
	}
	if s.URL != "https://env.example.com" || s.Username != "block-user" || s.Password != "env-pass" {
		t.Fatalf("unexpected settings %+v", s)
	}
	if s.Insecure {
		t.Fatal("TLS verification must be on by default")
	}
	t.Setenv("APPGATE_INSECURE", "true")
	if s, _ := resolveSettings(providerModel{}); !s.Insecure {
		t.Fatal("APPGATE_INSECURE=true must disable TLS verification")
	}
	if s, _ := resolveSettings(providerModel{Insecure: types.BoolValue(false)}); s.Insecure {
		t.Fatal("the provider block takes precedence over APPGATE_INSECURE")
	}
	t.Setenv("APPGATE_INSECURE", "")
	if s.apiVersion != "v23" || s.LoginTimeout != 10*time.Minute {
		t.Fatalf("client_version must map to api_version, got %q %s", s.apiVersion, s.LoginTimeout)
	}
	if _, err := uuid.Parse(s.DeviceID); err != nil {
		t.Fatalf("default device id %q is not a UUID", s.DeviceID)
	}

	file := filepath.Join(t.TempDir(), "appgate.json")
	if err := os.WriteFile(file, []byte(`{"appgate_url":"https://file.example.com","appgate_api_version":"v24","appgate_login_timeout":"30s"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	s, diags = resolveSettings(providerModel{ConfigPath: types.StringValue(file), APIVersion: types.StringValue("v25")})
	if diags.HasError() {
		t.Fatal(diags)
	}
	if s.URL != "https://file.example.com" || s.apiVersion != "v24" || s.LoginTimeout != 30*time.Second {
		t.Fatalf("config file must take precedence, got %+v", s)
	}

	for name, m := range map[string]providerModel{
		"bad timeout":   {LoginTimeout: types.StringValue("soon")},
		"bad device id": {DeviceID: types.StringValue("not-a-uuid")},
		"missing file":  {ConfigPath: types.StringValue(filepath.Join(t.TempDir(), "missing.json"))},
	} {
		if _, diags := resolveSettings(m); !diags.HasError() {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestProviderConfigurationErrors(t *testing.T) {
	srv := newServer(t, "v25")
	for name, tc := range map[string]struct{ config, want string }{
		"missing url": {`provider "appgatesdp" {
  url      = ""
  username = "a"
  password = "b"
  api_version = "v25"
}`, `controller URL is required`},
		"missing credentials": {fmt.Sprintf(`provider "appgatesdp" {
  url         = %q
  username    = ""
  password    = ""
  api_version = "v25"
}`, srv.URL), `username and password are required`},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("APPGATE_USERNAME", "")
			t.Setenv("APPGATE_PASSWORD", "")
			t.Setenv("APPGATE_ADDRESS", "")
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: protoV6Factories,
				Steps: []resource.TestStep{{
					Config:      tc.config + "\ndata \"appgatesdp_global_settings\" \"x\" {}\n",
					ExpectError: regexp.MustCompile(tc.want),
				}},
			})
		})
	}
}

// TestBearerToken: an externally obtained token is used as is.
func TestBearerToken(t *testing.T) {
	srv := newServer(t, "v25")
	token := srv.Token()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
provider "appgatesdp" {
  url          = %q
  bearer_token = %q
  api_version  = "v25"
}
data "appgatesdp_global_settings" "x" {}
`, srv.URL, token),
			Check: func(*terraform.State) error {
				for _, r := range srv.Requests() {
					if r.Path == "/login" {
						return fmt.Errorf("bearer token must not log in")
					}
				}
				return nil
			},
		}},
	})
}

func TestApplianceDeleteDeactivates(t *testing.T) {
	api, _ := generated.API("v25")
	srv := sdptest.New(t, api)
	id := srv.Seed("/appliances", map[string]any{"name": "gw", "activated": true})
	c := newTestClient(t, srv, 25)
	if err := deactivateAppliance(context.Background(), c, api.Resources["appliance"], id); err != nil {
		t.Fatal(err)
	}
	if srv.Object("/appliances", id)["activated"] != false {
		t.Fatal("appliance was not deactivated")
	}
	if lastRequest(srv, "POST", "/appliances/"+id+"/deactivate") == nil {
		var found bool
		for _, r := range srv.Requests() {
			found = found || (r.Path == "/appliances/"+id+"/deactivate" && r.Query == "wipe=true")
		}
		if !found {
			t.Fatal("deactivate must wipe the appliance")
		}
	}
	// Missing appliances and inactive appliances need no deactivation.
	if err := deactivateAppliance(context.Background(), c, api.Resources["appliance"], "missing"); err != nil {
		t.Fatal(err)
	}
}

func TestApplianceSeed(t *testing.T) {
	srv := newServer(t, "v25")
	inactive := srv.Seed("/appliances", map[string]any{"name": "gw", "activated": false})
	active := srv.Seed("/appliances", map[string]any{"name": "ctrl", "activated": true})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "v25") + fmt.Sprintf(`
data "appgatesdp_appliance_seed" "new" {
  appliance_id  = %q
  password      = "cz"
  validity_days = 2
}
data "appgatesdp_appliance_seed" "active" {
  appliance_id = %q
}
`, inactive, active),
			Check: resource.ComposeAggregateTestCheckFunc(
				// base64 of {"applianceId":"...","sshConfig":{"password":"cz","validityDays":2}}
				resource.TestCheckResourceAttrWith("data.appgatesdp_appliance_seed.new", "seed_file", func(v string) error {
					b, err := base64.StdEncoding.DecodeString(v)
					if err != nil {
						return err
					}
					if !strings.Contains(string(b), `"validityDays":2`) || !strings.Contains(string(b), inactive) {
						return fmt.Errorf("unexpected seed %q", b)
					}
					return nil
				}),
				resource.TestCheckResourceAttr("data.appgatesdp_appliance_seed.active", "seed_file", ""),
			),
		}},
	})
}

// TestEveryDataSourceReads reads every data source of every version: lookups
// by id of a seeded object, singletons, and read-only endpoints with their
// path parameters set.
func TestEveryDataSourceReads(t *testing.T) {
	forEachVersion(t, func(t *testing.T, v apiversion.Version) {
		api, _ := generated.API(v.APIVersion)
		for _, name := range sortedNames(api.DataSources) {
			m := api.DataSources[name]
			t.Run(name, func(t *testing.T) {
				parallelTerraform(t)
				srv := newServer(t, v.APIVersion)
				var args string
				switch {
				case m.ReadOnly:
					attrs := generatedDataSourceSchemas[name](context.Background()).Attributes
					for _, p := range m.Params {
						if !p.Required {
							continue
						}
						value := "x"
						if len(p.Enum) > 0 {
							value = p.Enum[0]
						}
						if _, isString := attrs[p.TF].(dschema.StringAttribute); isString {
							args += fmt.Sprintf("  %s = %q\n", p.TF, value)
						} else {
							args += fmt.Sprintf("  %s = 1\n", p.TF)
						}
					}
				case m.Singleton:
				default:
					obj := map[string]any{"name": "seeded"}
					if m.Discriminator != nil {
						obj[m.Discriminator.Property] = m.Discriminator.Value
					}
					id := uuid.NewString()
					obj["id"], obj[m.ID()] = id, id
					srv.Seed(m.Path, obj)
					args = fmt.Sprintf("  id = %q\n", id)
				}
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: protoV6Factories,
					Steps: []resource.TestStep{{
						Config: providerConfig(srv, v.APIVersion) + fmt.Sprintf("data %q \"test\" {\n%s}\n", "appgatesdp_"+name, args),
						Check:  resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttrSet("data.appgatesdp_"+name+".test", "id"), checkAcceptHeaders(srv, v)),
					}},
				})
			})
		}
	})
}

// TestTLSVerifiedByDefault: certificates are verified unless insecure is set;
// a private CA is trusted through pem_filepath.
func TestTLSVerifiedByDefault(t *testing.T) {
	api, _ := generated.API("v25")
	srv := sdptest.NewTLS(t, api)
	caFile := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(caFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := func(extra string) string {
		return fmt.Sprintf(`
provider "appgatesdp" {
  url         = %q
  username    = %q
  password    = %q
  api_version = "v25"
  %s
}
data "appgatesdp_global_settings" "test" {}
`, srv.URL, sdptest.Username, sdptest.Password, extra)
	}
	t.Setenv("APPGATE_INSECURE", "")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6Factories,
		Steps: []resource.TestStep{
			{Config: cfg(""), ExpectError: regexp.MustCompile(`(?s)certificate.*set\s+pem_filepath`)},
			{Config: cfg(fmt.Sprintf("pem_filepath = %q", caFile))},
			{Config: cfg("insecure = true")},
		},
	})
}
