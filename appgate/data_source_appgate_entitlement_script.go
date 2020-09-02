package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceEntitlementScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateEntitlementScriptRead,
		Schema: map[string]*schema.Schema{
			"entitlement script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement script_name"},
			},
			"entitlement script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"entitlement script_id"},
			},
		},
	}
}

func dataSourceAppgateEntitlementScriptRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.EntitlementScriptsApi

	entitlementScriptID, iok := d.GetOk("entitlement script_id")
	entitlementScriptName, nok := d.GetOk("entitlement script_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of entitlement script_id or entitlement script_name attributes")
	}
	var reqErr error
	var entitlementScript *openapi.EntitlementScript
	if iok {
		entitlementScript, reqErr = findEntitlementScriptByUUID(api, entitlementScriptID.(string), token)
	} else {
		entitlementScript, reqErr = findEntitlementScriptByName(api, entitlementScriptName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(entitlementScript.Id)
	d.Set("name", entitlementScript.Name)

	return nil
}

func findEntitlementScriptByUUID(api *openapi.EntitlementScriptsApiService, id string, token string) (*openapi.EntitlementScript, error) {
	entitlementScript, _, err := api.EntitlementScriptsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &entitlementScript, nil
}

func findEntitlementScriptByName(api *openapi.EntitlementScriptsApiService, name string, token string) (*openapi.EntitlementScript, error) {
	request := api.EntitlementScriptsGet(context.Background())

	entitlementScript, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range entitlementScript.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find Critera script %s", name)
}
