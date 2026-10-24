package acctest

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion/generated"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdp"
	"github.com/appgate/terraform-provider-appgatesdp/internal/sdptest"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestConformance(t *testing.T) {
	fields := map[string]*apimeta.Field{
		"name":  {JSON: "name", Type: "string"},
		"mode":  {JSON: "mode", Type: "string", Enum: []string{"a", "b"}},
		"port":  {JSON: "port", Type: "int64"},
		"tags":  {JSON: "tags", Type: "list", Elem: "string"},
		"rules": {JSON: "rules", Type: "list", Fields: map[string]*apimeta.Field{"on": {JSON: "on", Type: "bool"}}},
	}
	got := Summarize(Conformance("condition", decode(t, `{"name":"n","mode":"c","port":1.5,"tags":["x",1],"rules":[{"on":"yes","extra":1}],"undocumented":true}`), fields))
	for _, want := range []string{
		`warning condition.mode: value "c" not in enum [a b]`,
		"error condition.port: expected integer, got 1.5",
		"error condition.tags[1]: expected string",
		"error condition.rules[0].on: expected bool",
		"info condition.rules[0].extra: field not in specification",
		"info condition.undocumented: field not in specification",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if f := Conformance("ok", decode(t, `{"name":"n","mode":"a","port":2,"tags":[],"rules":[{"on":true}]}`), fields); len(f) != 0 {
		t.Errorf("unexpected findings %v", f)
	}
}

func TestSweepOnlyDeletesTestObjects(t *testing.T) {
	api, err := generated.API("v25")
	if err != nil {
		t.Fatal(err)
	}
	srv := sdptest.New(t, api)
	keep := srv.Seed("/conditions", map[string]any{"name": "production"})
	gone := srv.Seed("/conditions", map[string]any{"name": NamePrefix + "condition-1"})
	srv.Seed("/blacklist", map[string]any{"id": "CN=" + NamePrefix + "user,OU=local", "userDistinguishedName": "CN=" + NamePrefix + "user,OU=local"})
	realUser := srv.Seed("/blacklist", map[string]any{"id": "CN=bob,OU=ldap", "userDistinguishedName": "CN=bob,OU=ldap"})
	c, err := sdp.New(sdp.Config{URL: srv.URL, Username: sdptest.Username, Password: sdptest.Password}, 25)
	if err != nil {
		t.Fatal(err)
	}
	if err := Sweep(context.Background(), c, api); err != nil {
		t.Fatal(err)
	}
	if srv.Count("/blacklist") != 1 || srv.Object("/blacklist", realUser) == nil {
		t.Fatal("sweep must delete blacklist entries of test users only")
	}
	if srv.Object("/conditions", keep) == nil || srv.Object("/conditions", gone) != nil {
		t.Fatal("sweep must delete exactly the objects named with the test prefix")
	}
}

func TestEnvironments(t *testing.T) {
	t.Setenv("SDP_ACC_V24_URL", "https://v24.example.com")
	t.Setenv("SDP_ACC_V24_USERNAME", "admin")
	envs := Environments()
	if len(envs) != 1 || envs[0].Version.APIVersion != "v24" || envs[0].Provider != "local" {
		t.Fatalf("got %+v", envs)
	}
	if cfg := envs[0].ProviderConfig(); !strings.Contains(cfg, `api_version = "v24"`) {
		t.Fatalf("provider config %s", cfg)
	}
	if !strings.HasPrefix(Name("condition"), NamePrefix+"condition-") {
		t.Fatal("names must carry the sweep prefix")
	}
}
