package appgate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestReadAWSResolversFromConfigBooleans(t *testing.T) {
	names := resourceAppgateSite().Schema["name_resolution"].Elem.(*schema.Resource)
	fields := names.Schema["aws_resolvers"].Elem.(*schema.Resource).Schema
	for _, discovery := range []bool{false, true} {
		for _, iam := range []bool{false, true} {
			t.Run(fmt.Sprintf("discovery=%t/iam=%t", discovery, iam), func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, fields, map[string]interface{}{
					"name": "test", "vpc_auto_discovery": discovery, "use_iam_role": iam,
					"vpcs": []interface{}{"vpc-test"}, "regions": []interface{}{"us-east-1"},
				})
				raw := make(map[string]interface{})
				for key := range fields {
					raw[key] = d.Get(key)
				}
				rows, err := readAWSResolversFromConfig(Appliance65Version, []interface{}{raw})
				if err != nil {
					t.Fatal(err)
				}
				encoded, err := json.Marshal(rows[0])
				if err != nil {
					t.Fatal(err)
				}
				var payload map[string]interface{}
				if err := json.Unmarshal(encoded, &payload); err != nil {
					t.Fatal(err)
				}
				for key, want := range map[string]bool{"vpcAutoDiscovery": discovery, "useIAMRole": iam} {
					if got, present := payload[key]; !present || got != want {
						t.Errorf("%s = %v (present=%t), want %t; JSON: %s", key, got, present, want, encoded)
					}
				}
			})
		}
	}
	for _, nilValues := range []bool{false, true} {
		t.Run(fmt.Sprintf("unset/nil=%t", nilValues), func(t *testing.T) {
			raw := map[string]interface{}{"vpcs": []interface{}{}, "regions": []interface{}{}}
			if nilValues {
				raw["vpc_auto_discovery"] = nil
				raw["use_iam_role"] = nil
			}
			rows, err := readAWSResolversFromConfig(Appliance65Version, []interface{}{raw})
			if err != nil {
				t.Fatal(err)
			}
			if rows[0].HasVpcAutoDiscovery() || rows[0].HasUseIAMRole() {
				t.Fatal("absent/nil flags must remain unset")
			}
		})
	}
}

func TestAWSResolverBooleansSiteLifecycle(t *testing.T) {
	api, _, mux, _, _, teardown := setup()
	defer teardown()
	client := &Client{API: api, Config: &Config{BearerToken: "test-token"}, ApplianceVersion: Appliance65Version}
	var stored map[string]interface{}
	var writes []map[string]interface{}
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
			writes = append(writes, body)
			stored = body
			stored["id"] = "test-site"
			stored["vpn"] = map[string]interface{}{"snat": false}
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
	for _, tc := range []struct {
		name      string
		flags     map[string]interface{}
		discovery bool
		iam       bool
	}{
		{"create omitted", nil, false, false},
		{"enable both", map[string]interface{}{"vpc_auto_discovery": true, "use_iam_role": true}, true, true},
		{"disable discovery", map[string]interface{}{"vpc_auto_discovery": false, "use_iam_role": true}, false, true},
		{"disable IAM", map[string]interface{}{"vpc_auto_discovery": true, "use_iam_role": false}, true, false},
		{"disable both", map[string]interface{}{"vpc_auto_discovery": false, "use_iam_role": false}, false, false},
		{"re-enable both", map[string]interface{}{"vpc_auto_discovery": true, "use_iam_role": true}, true, true},
		{"omit after true", nil, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resolver := map[string]interface{}{
				"name": "test", "update_interval": 60,
				"vpcs": []interface{}{"vpc-test"}, "regions": []interface{}{"us-east-1"},
			}
			for key, value := range tc.flags {
				resolver[key] = value
			}
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "test", "name_resolution": []interface{}{map[string]interface{}{
					"aws_resolvers": []interface{}{resolver},
				}},
			})
			diff, err := resource.Diff(ctx, state, config, client)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil || diff.Empty() {
				t.Fatal("expected create/update diff")
			}
			previousWrites := len(writes)
			applied, diagnostics := resource.Apply(ctx, state, diff, client)
			if diagnostics.HasError() {
				t.Fatalf("apply failed: %v", diagnostics)
			}
			if len(writes) != previousWrites+1 {
				t.Fatal("expected exactly one POST/PUT")
			}
			payload := stored["nameResolution"].(map[string]interface{})["awsResolvers"].([]interface{})[0].(map[string]interface{})
			for key, want := range map[string]bool{"vpcAutoDiscovery": tc.discovery, "useIAMRole": tc.iam} {
				if got, present := payload[key]; !present || got != want {
					t.Errorf("request %s = %v (present=%t), want %t", key, got, present, want)
				}
			}
			refreshed, diagnostics := resource.RefreshWithoutUpgrade(ctx, applied, client)
			if diagnostics.HasError() {
				t.Fatalf("refresh failed: %v", diagnostics)
			}
			diff, err = resource.Diff(ctx, refreshed, config, client)
			if err != nil {
				t.Fatal(err)
			}
			if diff != nil && !diff.Empty() {
				t.Fatalf("unexpected drift after refresh: %#v", diff.Attributes)
			}
			// Each next step must start from the refreshed SDK state.
			state = refreshed
		})
	}
}
