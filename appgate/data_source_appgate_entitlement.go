package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateEntitlement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateEntitlementRead,
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

func dataSourceAppgateEntitlementRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementsApi

	entitlementID, iok := d.GetOk("entitlement_id")
	entitlementName, nok := d.GetOk("entitlement_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of entitlement_id or entitlement_name attributes")
	}
	var reqErr error
	var entitlement *openapi.Entitlement
	if iok {
		entitlement, reqErr = findEntitlementByUUID(api, entitlementID.(string), token)
	} else {
		entitlement, reqErr = findEntitlementByName(api, entitlementName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(entitlement.Id)
	d.Set("name", entitlement.Name)

	return nil
}

func findEntitlementByUUID(api *openapi.EntitlementsApiService, id string, token string) (*openapi.Entitlement, error) {
	entitlement, _, err := api.EntitlementsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &entitlement, nil
}

func findEntitlementByName(api *openapi.EntitlementsApiService, name string, token string) (*openapi.Entitlement, error) {
	request := api.EntitlementsGet(context.Background())

	entitlement, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range entitlement.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find entitlement %s", name)
}
