package appgate

import (
	"context"
	"reflect"
	"testing"
)

func TestResourceExampleInstanceStateUpgradeV0(t *testing.T) {
	testCases := []struct {
		Description   string
		InputState    map[string]interface{}
		ExpectedState map[string]interface{}
		Meta          interface{}
	}{
		{
			Description:   "missing state",
			InputState:    nil,
			ExpectedState: nil,
		},
		{
			Description: "move device_limit_per_user to root level",
			InputState: map[string]interface{}{
				"name":         "foobar",
				"notes":        "Managed by terraform",
				"object_class": "user",
				"on_boarding_two_factor": map[string]interface{}{
					"always_required":       false,
					"claim_suffix":          "onBoarding",
					"device_limit_per_user": 6,
					"message":               "welcome",
					"mfa_provider_id":       "3ae98d53-c520-437f-99e4-451f936e6d2c",
				},
			},
			ExpectedState: map[string]interface{}{
				"name":                  "foobar",
				"notes":                 "Managed by terraform",
				"object_class":          "user",
				"device_limit_per_user": 6,
				"on_boarding_two_factor": map[string]interface{}{
					"always_required": false,
					"claim_suffix":    "onBoarding",
					"message":         "welcome",
					"mfa_provider_id": "3ae98d53-c520-437f-99e4-451f936e6d2c",
				},
			},
			Meta: &Client{
				ApplianceVersion: Appliance64Version,
			},
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.Description, func(t *testing.T) {
			got, err := resourceIdentityProvidereUpgradeV0(context.Background(), testCase.InputState, testCase.Meta)
			if err != nil {
				t.Fatalf("error migrating state: %s", err)
			}
			if !reflect.DeepEqual(testCase.ExpectedState, got) {
				t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", testCase.ExpectedState, got)
			}
		})
	}
}
