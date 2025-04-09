package appgate

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeviceScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateDeviceScriptRead,
		Schema: map[string]*schema.Schema{
			"device_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"device_script_name"},
			},
			"device_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"device_script_id"},
			},
		},
	}
}

func dataSourceAppgateDeviceScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi
	deviceScript, diags := ResolveDeviceScriptFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(deviceScript.GetId())
	d.Set("device_script_name", deviceScript.GetName())
	d.Set("device_script_id", deviceScript.GetId())

	return nil
}
