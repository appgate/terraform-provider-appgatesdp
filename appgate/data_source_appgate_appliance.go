package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateAppliance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateApplianceRead,
		Schema: map[string]*schema.Schema{
			"appliance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"appliance_name"},
			},
			"appliance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"appliance_id"},
			},
		},
	}
}

func dataSourceAppgateApplianceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source Appliance")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	appliance, diags := ResolveApplianceFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(appliance.GetId())
	d.Set("appliance_name", appliance.GetName())
	d.Set("appliance_id", appliance.GetId())
	return nil
}
