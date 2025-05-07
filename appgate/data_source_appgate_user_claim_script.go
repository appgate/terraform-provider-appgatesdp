package appgate

import (
	"context"
	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserClaimScript() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateUserClaimScriptRead,
		Schema: map[string]*schema.Schema{
			"user_claim_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_claim_script_name"},
			},
			"user_claim_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_claim_script_id"},
			},
		},
	}
}

func dataSourceAppgateUserClaimScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.UserClaimScriptsApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	userClaimScript, diags := ResolveUserScriptFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(userClaimScript.GetId())
	d.Set("user_claim_script_name", userClaimScript.GetName())
	d.Set("user_claim_script_id", userClaimScript.GetId())

	return nil
}
