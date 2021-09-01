package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateEntitlement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateEntitlementRead,
		Schema: map[string]*schema.Schema{
			"entitlement_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_name"},
			},
			"entitlement_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement_id"},
			},
		},
	}
}

func dataSourceAppgateEntitlementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.EntitlementsApi

	entitlementID, iok := d.GetOk("entitlement_id")
	entitlementName, nok := d.GetOk("entitlement_name")

	if !iok && !nok {
		return diag.FromErr(fmt.Errorf("please provide one of entitlement_id or entitlement_name attributes"))
	}
	var reqErr error
	var entitlement *openapi.Entitlement
	if iok {
		entitlement, reqErr = findEntitlementByUUID(ctx, api, entitlementID.(string), token)
	} else {
		entitlement, reqErr = findEntitlementByName(ctx, api, entitlementName.(string), token)
	}
	if reqErr != nil {
		return diag.FromErr(reqErr)
	}
	d.SetId(entitlement.Id)
	d.Set("entitlement_name", entitlement.Name)
	d.Set("entitlement_id", entitlement.Id)

	return diags
}

func findEntitlementByUUID(ctx context.Context, api *openapi.EntitlementsApiService, id, token string) (*openapi.Entitlement, error) {
	entitlement, _, err := api.EntitlementsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func findEntitlementByName(ctx context.Context, api *openapi.EntitlementsApiService, name, token string) (*openapi.Entitlement, error) {
	request := api.EntitlementsGet(ctx).Query(name).Range_("0-50").OrderBy("name").Authorization(token)

	entitlements, _, err := request.Execute()
	if err != nil {
		return nil, err
	}

	for _, entitlement := range entitlements.GetData() {
		if entitlement.GetName() == name {
			return &entitlement, nil
		}
	}
	return nil, fmt.Errorf("Failed to find entitlement by name %q", name)
}
