package appgate

import (
	"context"
	"strings"

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

	api := meta.(*Client).API.ClientProfilesApi
	clientProfile, diags := ResolveClientProfileFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(clientProfile.GetId())
	d.Set("client_profile_name", clientProfile.GetName())
	d.Set("client_profile_id", clientProfile.GetId())

	if currentVersion.LessThan(Appliance61Version) {
		clientConnections, _, err := api.ClientConnectionsGet(ctx).Authorization(token).Execute()
		if err != nil {
			diags = AppendFromErr(diags, err)
			return diags
		}
		for _, profile := range clientConnections.GetProfiles() {
			if strings.EqualFold(profile.GetName(), clientProfile.GetName()) {
				d.Set("url", profile.GetUrl())
			}
		}

	}
	url, _, err := api.ClientProfilesIdUrlGet(ctx, clientProfile.GetId()).Authorization(token).Execute()
	if err != nil {
		diags = AppendFromErr(diags, err)
		return diags
	}
	d.Set("url", url.GetUrl())

	return nil
}
