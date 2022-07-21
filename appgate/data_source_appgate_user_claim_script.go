package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserClaimScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateUserClaimScriptRead,
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

func dataSourceAppgateUserClaimScriptRead(d *schema.ResourceData, meta interface{}) error {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.UserClaimScriptsApi

	userClaimScriptID, iok := d.GetOk("user_claim_script_id")
	userClaimScriptName, nok := d.GetOk("user_claim_script_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of user_claim_script_id or user_claim_script_name attributes")
	}
	var reqErr error
	var userClaimScript *openapi.UserScript
	if iok {
		userClaimScript, reqErr = findUserClaimScriptByUUID(api, userClaimScriptID.(string), token)
	} else {
		userClaimScript, reqErr = findUserClaimScriptByName(api, userClaimScriptName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(userClaimScript.GetId())
	d.Set("user_claim_script_name", userClaimScript.GetName())

	return nil
}

func findUserClaimScriptByUUID(api *openapi.UserClaimScriptsApiService, id string, token string) (*openapi.UserScript, error) {
	userClaimScript, _, err := api.UserScriptsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return userClaimScript, nil
}

func findUserClaimScriptByName(api *openapi.UserClaimScriptsApiService, name string, token string) (*openapi.UserScript, error) {
	request := api.UserScriptsGet(context.Background())

	userClaimScript, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range userClaimScript.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find User claim script %s", name)
}
