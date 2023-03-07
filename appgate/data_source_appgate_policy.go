package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgatePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgatePolicyRead,
		Schema: map[string]*schema.Schema{

			"policy_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"policy_name"},
			},

			"policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"policy_id"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAppgatePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source policy")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.PoliciesApi
	policy, diags := ResolvePolicyFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(policy.GetId())
	d.Set("name", policy.GetName())
	d.Set("policy_id", policy.GetId())
	d.Set("tags", policy.GetTags())

	return nil
}
