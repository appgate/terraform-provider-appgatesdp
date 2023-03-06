package appgate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCriteriaScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateCriteriaScriptRead,
		Schema: map[string]*schema.Schema{
			"criteria_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"criteria_script_name"},
			},
			"criteria_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"criteria_script_id"},
			},
		},
	}
}

func dataSourceAppgateCriteriaScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.CriteriaScriptsApi
	criteraScript, diags := ResolveCriteriaScriptFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(criteraScript.GetId())
	d.Set("criteria_script_name", criteraScript.GetName())
	d.Set("criteria_script_id", criteraScript.GetId())

	return nil
}
