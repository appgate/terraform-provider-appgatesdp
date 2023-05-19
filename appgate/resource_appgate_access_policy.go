package appgate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateAccessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return resourceAppgatePolicyCreate(context.WithValue(ctx, PolicyTypeCtx, PolicyTypeAccess), rd, meta)
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
				basePolicyEntitlementAttributes(),
				basePolicyDeploymentSiteAttributes(),
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
