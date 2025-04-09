package appgate

import (
	"context"
	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateCondition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateConditionRead,
		Schema: map[string]*schema.Schema{
			"condition_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"condition_name"},
			},
			"condition_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"condition_id"},
			},
		},
	}
}

func dataSourceAppgateConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ConditionsApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	condition, diags := ResolveConditionFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}
	d.SetId(condition.GetId())
	d.Set("condition_id", condition.GetId())
	d.Set("condition_name", condition.GetName())

	return nil
}
