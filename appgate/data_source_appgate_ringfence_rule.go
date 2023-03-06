package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateRingfenceRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateRingfenceRuleRead,
		Schema: map[string]*schema.Schema{

			"ringfence_rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ringfence_rule_name"},
			},

			"ringfence_rule_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ringfence_rule_id"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAppgateRingfenceRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source Ringfence Rules")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.RingfenceRulesApi
	ringfenceRule, diags := ResolveRingfenceRuleFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(ringfenceRule.GetId())
	d.Set("ringfence_rule_id", ringfenceRule.GetId())
	d.Set("ringfence_rule_name", ringfenceRule.GetName())
	d.Set("tags", ringfenceRule.GetTags())

	return nil
}
