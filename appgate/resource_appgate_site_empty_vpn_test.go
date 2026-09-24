package appgate

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
