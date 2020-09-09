package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCriteriaScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateCriteriaScriptRead,
		Schema: map[string]*schema.Schema{
			"criteria_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"criteria_script_name"},
			},
			"criteria_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"criteria_script_id"},
			},
		},
	}
}

func dataSourceAppgateCriteriaScriptRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.CriteriaScriptsApi

	criteraScriptID, iok := d.GetOk("criteria_script_id")
	criteraScriptName, nok := d.GetOk("criteria_script_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of criteria_script_id or criteria_script_name attributes")
	}
	var reqErr error
	var criteraScript *openapi.CriteriaScript
	if iok {
		criteraScript, reqErr = findCriteriaScriptByUUID(api, criteraScriptID.(string), token)
	} else {
		criteraScript, reqErr = findCriteriaScriptByName(api, criteraScriptName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(criteraScript.Id)
	d.Set("name", criteraScript.Name)

	return nil
}

func findCriteriaScriptByUUID(api *openapi.CriteriaScriptsApiService, id string, token string) (*openapi.CriteriaScript, error) {
	criteraScript, _, err := api.CriteriaScriptsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &criteraScript, nil
}

func findCriteriaScriptByName(api *openapi.CriteriaScriptsApiService, name string, token string) (*openapi.CriteriaScript, error) {
	request := api.CriteriaScriptsGet(context.Background())

	criteraScript, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range criteraScript.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find Critera script %s", name)
}
