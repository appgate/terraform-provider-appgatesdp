package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func identityProviderResourcev0() *schema.Resource {
	s := ldapProviderSchema()
	return &schema.Resource{
		Schema: s,
	}
}

func resourceIdentityProvidereUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return nil, nil
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if currentVersion.GreaterThanOrEqual(Appliance55Version) {
		if v, ok := rawState["on_boarding_two_factor"]; ok {
			twoFA := v.(map[string]interface{})
			if v, ok := twoFA["device_limit_per_user"]; ok {
				log.Printf("[INFO] Migrating device_limit_per_user to root level from on_boarding_two_factor")
				rawState["device_limit_per_user"] = v.(int)
				delete(twoFA, "device_limit_per_user")
				rawState["on_boarding_two_factor"] = twoFA
			}
		}
	}

	return rawState, nil
}
