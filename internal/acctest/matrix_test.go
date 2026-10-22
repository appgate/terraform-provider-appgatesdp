package acctest_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"golang.org/x/crypto/ssh"

	"github.com/appgate/terraform-provider-appgatesdp/internal/acctest"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
)

// selfSignedPEM returns a throwaway CA certificate for certificate fields.
func selfSignedPEM(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()), Subject: pkix.Name{CommonName: acctest.NamePrefix + "ca"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * time.Hour),
		IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})))
}

// brokeredConfig is a brokered entitlement on the default site with a
// condition, which entitlements require.
func brokeredConfig(kind string, port int) string {
	return fmt.Sprintf(`data "appgatesdp_site" "default" {
  name = "Default Site"
}
resource "appgatesdp_condition" "always" {
  name       = "%%[1]s-condition"
  expression = "return true;"
}
resource "appgatesdp_brokered_%[1]s_entitlement" "test" {
  name         = %%[1]q
  site         = data.appgatesdp_site.default.id
  conditions   = [appgatesdp_condition.always.id]
  host         = "%[1]s.invalid"
  port         = %[2]d
  app_shortcut = { name = %%[1]q }
}`, kind, port)
}

// sshKeyPEM returns a throwaway OpenSSH private key.
func sshKeyPEM(t *testing.T) string {
	t.Helper()
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKey(key, "")
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(pem.EncodeToMemory(block)))
}

type matrixCase struct {
	// config is a resource block named "test"; %[1]s is a unique name,
	// %[2]s a PEM certificate and %[3]s an SSH private key.
	config string
	// ignore lists write-only attributes the API never returns.
	ignore []string
	// identifier uses a JavaScript identifier as name.
	identifier bool
}

var matrix = map[string]matrixCase{
	"criteria_script": {identifier: true, config: `resource "appgatesdp_criteria_script" "test" {
  name = %[1]q
  expression = "return true;"
}`},
	"entitlement_script": {config: `resource "appgatesdp_entitlement_script" "test" {
  name = %[1]q
  expression = "return [];"
}`},
	"user_claim_script": {config: `resource "appgatesdp_user_claim_script" "test" {
  name = %[1]q
  expression = "return {};"
}`},
	"device_script": {config: `resource "appgatesdp_device_script" "test" {
  name     = %[1]q
  filename = "tf-acc.sh"
  file     = "` + base64.StdEncoding.EncodeToString([]byte("#!/bin/sh\necho ok\n")) + `"
}`, ignore: []string{"file"}},
	"replication_target": {config: `resource "appgatesdp_replication_target" "test" {
  name             = %[1]q
  replication_tags = ["tf-acc"]
}`},
	"service_user": {config: `resource "appgatesdp_service_user" "test" {
  name     = %[1]q
  password = "Tf-acc-Passw0rd!"
}`, ignore: []string{"password"}},
	"tag": {config: `resource "appgatesdp_tag" "test" {
  name = %[1]q
}`},
	"client_profile_group": {config: `resource "appgatesdp_client_profile_group" "test" {
  name = %[1]q
}`},
	"service_database_identity_provider": {config: `resource "appgatesdp_service_database_identity_provider" "test" {
  name = %[1]q
}`},
	"protocol_template": {config: `resource "appgatesdp_protocol_template" "test" {
  name    = %[1]q
  subtype = "tcp_up"
  ports   = ["443"]
}`},
	"credentials_secret": {config: `resource "appgatesdp_credentials_secret" "test" {
  name     = %[1]q
  username = "tf-acc"
  password = "tf-acc-secret"
}`, ignore: []string{"password"}},
	"generic_secret": {config: `resource "appgatesdp_generic_secret" "test" {
  name   = %[1]q
  secret = "tf-acc-secret"
}`, ignore: []string{"secret"}},
	"ssh_key_secret": {config: `resource "appgatesdp_ssh_key_secret" "test" {
  name        = %[1]q
  username    = "tf-acc"
  private_key = <<-KEY
%[3]s
KEY
}`, ignore: []string{"private_key"}},
	"brokered_rdp_entitlement": {config: brokeredConfig("rdp", 3389)},
	"brokered_ssh_entitlement": {config: brokeredConfig("ssh", 22)},
	"site": {config: `resource "appgatesdp_site" "test" {
  name            = %[1]q
  network_subnets = ["10.99.0.0/24"]
}`},
	"ip_pool": {config: `resource "appgatesdp_ip_pool" "test" {
  name   = %[1]q
  ranges = [{ first = "10.250.0.10", last = "10.250.0.200" }]
}`},
	"ringfence_rule": {config: `resource "appgatesdp_ringfence_rule" "test" {
  name = %[1]q
  actions = [{
    protocol  = "tcp"
    direction = "out"
    action    = "allow"
    hosts     = ["10.0.0.0/16"]
    ports     = ["443"]
  }]
}`},
	"administrative_role": {config: `resource "appgatesdp_administrative_role" "test" {
  name       = %[1]q
  privileges = [{ type = "View", target = "Condition" }]
}`},
	"local_user": {config: `resource "appgatesdp_local_user" "test" {
  name       = %[1]q
  first_name = "Terraform"
  last_name  = "Acceptance"
  password   = "Tf-acc-Passw0rd!"
}`, ignore: []string{"password"}},
	"mfa_provider": {config: `resource "appgatesdp_mfa_provider" "test" {
  name          = %[1]q
  type          = "Radius"
  hostnames     = ["radius.invalid"]
  port          = 1812
  mode          = "OneFactor"
  shared_secret = "tf-acc-secret"
}`, ignore: []string{"shared_secret"}},
	"trusted_certificate": {config: `resource "appgatesdp_trusted_certificate" "test" {
  name = %[1]q
  pem  = <<-PEM
%[2]s
PEM
}`},
	"access_policy": {config: `resource "appgatesdp_access_policy" "test" {
  name = %[1]q
  expression = "return false;"
}`},
	"admin_policy": {config: `resource "appgatesdp_admin_policy" "test" {
  name = %[1]q
  expression = "return false;"
}`},
	"device_policy": {config: `resource "appgatesdp_device_policy" "test" {
  name = %[1]q
  expression = "return false;"
}`},
	"dns_policy": {config: `resource "appgatesdp_dns_policy" "test" {
  name = %[1]q
  expression = "return false;"
}`},
	"stop_policy": {config: `resource "appgatesdp_stop_policy" "test" {
  name = %[1]q
  expression = "return false;"
}`},
	"ldap_identity_provider": {config: `resource "appgatesdp_ldap_identity_provider" "test" {
  name                     = %[1]q
  hostnames                = ["ldap.invalid"]
  port                     = 389
  admin_distinguished_name = "CN=admin,DC=example,DC=com"
  admin_password           = "tf-acc-secret"
  base_dn                  = "DC=example,DC=com"
}`, ignore: []string{"admin_password"}},
	"radius_identity_provider": {config: `resource "appgatesdp_radius_identity_provider" "test" {
  name          = %[1]q
  hostnames     = ["radius.invalid"]
  shared_secret = "tf-acc-secret"
}`, ignore: []string{"shared_secret"}},
	"oidc_identity_provider": {config: `resource "appgatesdp_oidc_identity_provider" "test" {
  name     = %[1]q
  issuer   = "https://idp.invalid"
  audience = "tf-acc"
}`},
	"appliance": {config: `resource "appgatesdp_appliance" "test" {
  name     = %[1]q
  hostname = "tf-acc-appliance.invalid"
  client_interface = {
    hostname = "tf-acc-appliance.invalid"
  }
  networking = {
    nics = [{
      enabled = true
      name    = "eth0"
      ipv4    = { dhcp = { enabled = true } }
    }]
  }
}
data "appgatesdp_appliance_seed" "test" {
  appliance_id = appgatesdp_appliance.test.id
  password     = "tf-acc-cz-password"
}`, ignore: []string{"updated"}}, // exporting the seed updates the appliance
}

// TestAccResourceMatrix creates, re-plans (must be empty), imports and
// destroys a realistic minimal object of each resource kind.
func TestAccResourceMatrix(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		sweepAfter(t, env)
		_, api := env.Client(t)
		pemCert, sshKey := selfSignedPEM(t), sshKeyPEM(t)
		for _, name := range sortedNames(matrix) {
			tc := matrix[name]
			t.Run(name, func(t *testing.T) {
				if _, ok := api.Resources[name]; !ok {
					t.Skipf("not available in %s", env.Version.APIVersion)
				}
				addr := "appgatesdp_" + name + ".test"
				objName := acctest.Name(strings.ReplaceAll(name, "_", "-"))
				if tc.identifier {
					objName = acctest.Identifier(name)
				}
				cfg := env.ProviderConfig() + fmt.Sprintf(tc.config, objName, pemCert, sshKey) +
					// Every resource has a data source reading it back by id.
					fmt.Sprintf("\ndata %q \"by_id\" {\n  id = %s.id\n}\n", "appgatesdp_"+name, addr)
				check := resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(addr, "id"),
					resource.TestCheckResourceAttrPair("data.appgatesdp_"+name+".by_id", "id", addr, "id"),
				)
				if name == "appliance" {
					check = resource.ComposeAggregateTestCheckFunc(check, resource.TestCheckResourceAttrWith("data.appgatesdp_appliance_seed.test", "seed_file", func(v string) error {
						if v == "" {
							return fmt.Errorf("seed of an inactive appliance must not be empty")
						}
						return nil
					}))
				}
				resource.Test(t, resource.TestCase{
					ProtoV6ProviderFactories: factories,
					ErrorCheck:               acctest.SkipUnlicensed(t),
					Steps: []resource.TestStep{
						{Config: cfg, Check: check},
						{ResourceName: addr, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: tc.ignore},
					},
				})
			})
		}
	})
}

// TestAccBuiltInDataSources reads objects every collective has.
func TestAccBuiltInDataSources(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: factories,
			Steps: []resource.TestStep{{
				Config: env.ProviderConfig() + `
data "appgatesdp_identity_provider" "local" { name = "local" }
data "appgatesdp_global_settings" "this" {}
data "appgatesdp_certificate_authority" "ca" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.appgatesdp_identity_provider.local", "id"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_global_settings.this", "claims_token_expiration"),
					resource.TestCheckResourceAttrSet("data.appgatesdp_certificate_authority.ca", "subject"),
				),
			}},
		})
	})
}

// TestAccSingletons applies every settings singleton without attributes. The
// provider then PUTs the current settings back, so nothing changes; the plan
// must be empty afterwards and the singleton must import. Singletons whose
// destroy resets collective settings, and registrations, are left out.
func TestAccSingletons(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		acctest.RequireDestructive(t)
		_, api := env.Client(t)
		for _, name := range sortedNames(api.Resources) {
			r := api.Resources[name]
			if !r.Singleton || r.Delete || r.CreateMethod != "" {
				continue
			}
			t.Run(name, func(t *testing.T) {
				addr := "appgatesdp_" + name + ".test"
				resource.Test(t, resource.TestCase{
					ProtoV6ProviderFactories: factories,
					ErrorCheck:               acctest.SkipUnlicensed(t),
					Steps: []resource.TestStep{
						{Config: env.ProviderConfig() + fmt.Sprintf("resource %q \"test\" {}\ndata %q \"test\" {}\n", "appgatesdp_"+name, "appgatesdp_"+name)},
						{ResourceName: addr, ImportState: true, ImportStateId: name, ImportStateVerify: true},
					},
				})
			})
		}
	})
}

// TestAccReadOnlyDataSources reads every read-only endpoint data source of
// the version. Path parameters are taken from an existing object of the
// parent collection (e.g. a site for /sites/{id}/resources).
func TestAccReadOnlyDataSources(t *testing.T) {
	acctest.ForEachEnv(t, func(t *testing.T, env acctest.Env) {
		c, api := env.Client(t)
		for _, name := range sortedNames(api.DataSources) {
			ds := api.DataSources[name]
			if !ds.ReadOnly {
				continue
			}
			t.Run(name, func(t *testing.T) {
				var args string
				for _, p := range ds.Params {
					if !p.Required {
						continue
					}
					var value string
					switch {
					case p.In == "path":
						value = firstObjectID(t, c, ds.Path, p.Name)
					case len(p.Enum) > 0:
						value = p.Enum[0]
					default:
						t.Skipf("required parameter %s has no known value", p.Name)
					}
					args += fmt.Sprintf("  %s = %q\n", p.TF, value)
				}
				resource.Test(t, resource.TestCase{
					ProtoV6ProviderFactories: factories,
					ErrorCheck:               acctest.SkipUnlicensed(t),
					Steps: []resource.TestStep{{
						Config: env.ProviderConfig() + fmt.Sprintf("data %q \"test\" {\n%s}\n", "appgatesdp_"+name, args),
						Check:  resource.TestCheckResourceAttrSet("data.appgatesdp_"+name+".test", "id"),
					}},
				})
			})
		}
	})
}

// firstObjectID returns the id of the first object of the collection that
// precedes {param} in path, or skips the test.
func firstObjectID(t *testing.T, c *sdp.Client, path, param string) string {
	t.Helper()
	parent, _, _ := strings.Cut(path, "/{"+param+"}")
	if param != "id" || strings.Contains(parent, "{") {
		t.Skipf("no generic value for path parameter %s of %s", param, path)
	}
	var list struct {
		Data []map[string]any `json:"data"`
	}
	if err := c.Do(context.Background(), http.MethodGet, parent, url.Values{"range": {"0-1"}}, nil, &list); err != nil {
		t.Skipf("cannot list %s for a parameter value: %v", parent, err)
	}
	if len(list.Data) == 0 {
		t.Skipf("%s has no objects to read %s for", parent, path)
	}
	return fmt.Sprint(list.Data[0]["id"])
}
