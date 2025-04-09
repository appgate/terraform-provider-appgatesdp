package appgate

import (
	"context"
	"errors"
	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClientProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateClientProfileRead,
		Schema: map[string]*schema.Schema{
			"client_profile_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_profile_name"},
			},
			"client_profile_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_profile_id"},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateClientProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if currentVersion.LessThan(Appliance61Version) {
		return diag.FromErr(errors.New("data source appgatesdp_client_profile is not available on your version"))
	}

	api := meta.(*Client).API.ClientProfilesApi
	clientProfile, diags := ResolveClientProfileFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(clientProfile.GetId())
	d.Set("client_profile_name", clientProfile.GetName())
	d.Set("client_profile_id", clientProfile.GetId())

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	url, _, err := api.ClientProfilesIdUrlGet(ctx, clientProfile.GetId()).Execute()
	if err != nil {
		diags = AppendFromErr(diags, err)
		return diags
	}
	d.Set("url", url.GetUrl())

	return nil
}
