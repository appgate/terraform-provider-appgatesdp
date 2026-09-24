package appgate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
