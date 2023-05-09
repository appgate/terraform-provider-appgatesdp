package appgate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateDevicePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
			// we need to call GetToken to be able to get the Appliance version
			_, err := meta.(*Client).GetToken()
			if err != nil {
				return diag.FromErr(err)
			}
			currentVersion := meta.(*Client).ApplianceVersion
			if currentVersion.LessThan(Appliance55Version) {
				return diag.Errorf("appgatesdp_device_policy is not supported on your version")
			}

			return resourceAppgatePolicyCreate(context.WithValue(ctx, PolicyTypeCtx, PolicyTypeDevice), rd, meta)
		},
		ReadContext:   resourceAppgatePolicyRead,
		UpdateContext: resourceAppgatePolicyUpdate,
		DeleteContext: resourceAppgatePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: func() map[string]*schema.Schema {
			s := mergeSchemaMaps(
				basePolicySchema(),
				basePolicyClientAttributes(),
				basePolicyDeviceAttributes(),
				basePolicyRingfenceAttributes(),
			)
			s["expression"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  emptyPolicyExpression,
			}
			// Type is computed in CreateContext
			s["type"] = &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			}
			return s
		}(),
	}
}
