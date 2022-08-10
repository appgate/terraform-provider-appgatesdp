package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeviceScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateDeviceScriptRead,
		Schema: map[string]*schema.Schema{
			"device_script_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"device_script_name"},
			},
			"device_script_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"device_script_id"},
			},
		},
	}
}

func dataSourceAppgateDeviceScriptRead(d *schema.ResourceData, meta interface{}) error {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi

	deviceScriptID, iok := d.GetOk("device_script_id")
	deviceScriptName, nok := d.GetOk("device_script_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of device_script_id or device_script_name attributes")
	}
	var reqErr error
	var deviceScript *openapi.DeviceScript
	if iok {
		deviceScript, reqErr = findDeviceScriptByUUID(api, deviceScriptID.(string), token)
	} else {
		deviceScript, reqErr = findDeviceScriptByName(api, deviceScriptName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(deviceScript.GetId())
	d.Set("device_script_name", deviceScript.GetName())

	return nil
}

func findDeviceScriptByUUID(api *openapi.DeviceClaimScriptsApiService, id string, token string) (*openapi.DeviceScript, error) {
	deviceScript, _, err := api.DeviceScriptsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return deviceScript, nil
}

func findDeviceScriptByName(api *openapi.DeviceClaimScriptsApiService, name string, token string) (*openapi.DeviceScript, error) {
	request := api.DeviceScriptsGet(context.Background())

	deviceScript, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range deviceScript.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find Device script %s", name)
}
