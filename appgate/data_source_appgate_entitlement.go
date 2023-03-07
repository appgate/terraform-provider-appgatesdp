package appgate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateEntitlement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateEntitlementRead,
		Schema: map[string]*schema.Schema{
			"entitlement_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_name"},
			},
			"entitlement_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_id"},
			},
		},
	}
}

func dataSourceAppgateEntitlementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	entitlement, diags := ResolveEntitlementFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}
	d.SetId(entitlement.GetId())
	d.Set("entitlement_name", entitlement.GetName())
	d.Set("entitlement_id", entitlement.GetId())

	return diags
}
