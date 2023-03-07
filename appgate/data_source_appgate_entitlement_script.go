package appgate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEntitlementScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateEntitlementScriptRead,
		Schema: map[string]*schema.Schema{
			"entitlement_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_script_name"},
			},
			"entitlement_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_script_id"},
			},
		},
	}
}

func dataSourceAppgateEntitlementScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementScriptsApi
	entitlementScript, diags := ResolveEntitlementScriptFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(entitlementScript.GetId())
	d.Set("entitlement_script_id", entitlementScript.GetId())
	d.Set("entitlement_script_name", entitlementScript.GetName())

	return nil
}
