package appgate

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestReadSiteVPNFromConfigEmptyProtocols(t *testing.T) {
	for _, tc := range []struct {
		name      string
		protocols map[string]bool
	}{
		{"all omitted", map[string]bool{}},
		{"tls enabled without dtls", map[string]bool{"tls": true}},
		{"tls explicitly disabled", map[string]bool{"tls": false}},
		{"dtls explicitly disabled", map[string]bool{"dtls": false}},
		{"dtls enabled without tls", map[string]bool{"dtls": true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{"snat": true}
			for protocol, enabled := range tc.protocols {
				raw[protocol] = []interface{}{map[string]interface{}{"enabled": enabled}}
			}
			d := schema.TestResourceDataRaw(t, resourceAppgateSite().Schema, map[string]interface{}{
				"name": "test", "vpn": []interface{}{raw},
			})
			vpn, err := readSiteVPNFromConfig(d.Get("vpn").([]interface{}))
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(vpn)
			if err != nil {
				t.Fatal(err)
			}
			var payload map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &payload); err != nil {
				t.Fatal(err)
			}
			for _, protocol := range []string{"tls", "dtls"} {
				want, configured := tc.protocols[protocol]
				got, present := payload[protocol]
				if present != configured {
					t.Errorf("%s presence = %t, want %t: %s", protocol, present, configured, encoded)
				} else if configured {
					expected, _ := json.Marshal(map[string]bool{"enabled": want})
					if string(got) != string(expected) {
						t.Errorf("%s = %s, want %s", protocol, got, expected)
					}
				}
			}
			if string(payload["snat"]) != "true" {
				t.Errorf("SNAT changed: %s", encoded)
			}
		})
	}
}

func TestSiteVPNProtocolsLifecycle(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		name := "no DTLS in API"
		if legacy {
			name = "legacy API defaults"
		}
		t.Run(name, func(t *testing.T) {
			api, _, mux, _, _, teardown := setup()
			defer teardown()
			client := &Client{API: api, Config: &Config{BearerToken: "test-token"}, ApplianceVersion: Appliance65Version}
			var stored map[string]interface{}
			writes := 0
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPost, http.MethodPut:
					var body map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					vpn := body["vpn"].(map[string]interface{})
					_, dtls := vpn["dtls"]
					if !legacy && dtls {
						t.Error("DTLS was sent to an API that does not support it")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if writes == 0 {
						if _, tls := vpn["tls"]; tls || dtls {
							t.Error("create synthesized omitted protocol blocks")
						}
					} else {
						wantTLS := writes == 2
						if vpn["tls"].(map[string]interface{})["enabled"] != wantTLS {
							t.Errorf("update did not preserve explicit TLS %t", wantTLS)
						}
						if legacy && (!dtls || vpn["dtls"].(map[string]interface{})["enabled"] != false) {
							t.Error("update lost the legacy API's computed DTLS default")
						}
					}
					writes++
					stored = body
					stored["id"] = "test-site"
					stored["nameResolution"] = map[string]interface{}{"useHostsFile": false}
					// Model controller readback, including older API defaults.
					if _, present := vpn["tls"]; !present {
						vpn["tls"] = map[string]interface{}{"enabled": true}
					}
					if legacy && !dtls {
						vpn["dtls"] = map[string]interface{}{"enabled": false}
					}
				case http.MethodGet:
				default:
					t.Errorf("unexpected method: %s", r.Method)
				}
				if err := json.NewEncoder(w).Encode(stored); err != nil {
					t.Error(err)
				}
			}
			mux.HandleFunc("/sites", handler)
			mux.HandleFunc("/sites/test-site", handler)
			resource := resourceAppgateSite()
			ctx := context.Background()
			var state *terraform.InstanceState
			for step := 0; step < 3; step++ {
				vpn := map[string]interface{}{"snat": step == 0}
				if step > 0 {
					vpn["tls"] = []interface{}{map[string]interface{}{"enabled": step == 2}}
				}
				config := terraform.NewResourceConfigRaw(map[string]interface{}{
					"name": "test", "vpn": []interface{}{vpn},
				})
				diff, err := resource.Diff(ctx, state, config, client)
				if err != nil {
					t.Fatal(err)
				}
				if diff == nil || diff.Empty() {
					t.Fatal("expected create/update diff")
				}
				applied, diags := resource.Apply(ctx, state, diff, client)
				if diags.HasError() {
					t.Fatalf("apply failed: %v", diags)
				}
				state, diags = resource.RefreshWithoutUpgrade(ctx, applied, client)
				if diags.HasError() {
					t.Fatalf("refresh failed: %v", diags)
				}
				diff, err = resource.Diff(ctx, state, config, client)
				if err != nil {
					t.Fatal(err)
				}
				if diff != nil && !diff.Empty() {
					t.Fatalf("unexpected drift after refresh: %#v", diff.Attributes)
				}
			}
			if writes != 3 {
				t.Fatalf("got %d writes, want one create and two updates", writes)
			}
		})
	}
}
