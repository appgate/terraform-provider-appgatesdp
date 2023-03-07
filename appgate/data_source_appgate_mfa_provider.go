package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateMfaProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateMfaProviderRead,
		Schema: map[string]*schema.Schema{
			"mfa_provider_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mfa_provider_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateMfaProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source MFA provider")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.MFAProvidersApi
	provider, diags := ResolveMfaProviderFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(provider.GetId())
	d.Set("mfa_provider_name", provider.GetName())
	d.Set("mfa_provider_id", provider.GetId())
	return nil
}
