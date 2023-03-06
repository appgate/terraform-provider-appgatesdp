package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateApplianceCustomization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateApplianceCustomizationRead,
		Schema: map[string]*schema.Schema{
			"appliance_customization_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"appliance_customization_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateApplianceCustomizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source Appliance customization")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ApplianceCustomizationsApi
	appliance, diags := ResolveApplianceCustomizationFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(appliance.GetId())
	d.Set("appliance_customization_name", appliance.GetName())
	d.Set("appliance_customization_id", appliance.GetId())
	return nil
}
