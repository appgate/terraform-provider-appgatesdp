package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"

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

	adminID, iok := d.GetOk("administrative_role_id")
	adminName, nok := d.GetOk("administrative_role_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of administrative_role_id or administrative_role_name attributes")
	}
	var reqErr error
	var admin *openapi.AdministrativeRole
	if iok {
		admin, reqErr = findAdministrativeRoleByUUID(api, adminID.(string), token)
	} else {
		admin, reqErr = findAdministrativeRoleByName(api, adminName.(string), token)
	}
	if reqErr != nil {
		return diag.FromErr(reqErr)
	}
	log.Printf("[DEBUG] Got Administrative role: %+v", admin)

	d.SetId(admin.GetId())
	d.Set("administrative_role_name", admin.GetName())
	d.Set("administrative_role_id", admin.GetId())
	return nil
}

func findAdministrativeRoleByUUID(api *openapi.AdminRolesApiService, id string, token string) (*openapi.AdministrativeRole, error) {
	log.Printf("[DEBUG] Data source Administrative role get by UUID %s", id)
	admin, _, err := api.AdministrativeRolesIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func findAdministrativeRoleByName(api *openapi.AdminRolesApiService, name string, token string) (*openapi.AdministrativeRole, error) {
	log.Printf("[DEBUG] Data Administrative role get by name %s", name)
	request := api.AdministrativeRolesGet(context.Background())

	admin, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range admin.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find Administrative role %s", name)
}
