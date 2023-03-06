package appgate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateAdministrativeRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateAdministrativeRoleRead,
		Schema: map[string]*schema.Schema{
			"administrative_role_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"administrative_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateAdministrativeRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Data source Administrative role")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AdminRolesApi
	admin, diags := ResolveAdministrativeRoleFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(admin.GetId())
	d.Set("administrative_role_name", admin.GetName())
	d.Set("administrative_role_id", admin.GetId())
	return nil
}
