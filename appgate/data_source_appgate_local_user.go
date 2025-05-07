package appgate

import (
	"context"
	"github.com/appgate/sdp-api-client-go/api/v22/openapi"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateLocalUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateLocalUserRead,
		Schema: map[string]*schema.Schema{
			"local_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateLocalUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source local user")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LocalUsersApi
	localUser, diags := ResolveLocalUserFromResourceData(context.WithValue(ctx, openapi.ContextAccessToken, token), d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(localUser.GetId())
	d.Set("local_user_name", localUser.GetName())
	d.Set("local_user_id", localUser.GetId())
	return nil
}
