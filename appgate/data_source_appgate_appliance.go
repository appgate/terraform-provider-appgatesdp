package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAppgateAppliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateApplianceRead,
		Schema: map[string]*schema.Schema{
			"appliance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"appliance_name"},
			},
			"appliance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"appliance_id"},
			},
		},
	}
}

func dataSourceAppgateApplianceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source Appliance")
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi

	applianceID, iok := d.GetOk("appliance_id")
	applianceName, nok := d.GetOk("appliance_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of appliance_id or appliance_name attributes")
	}
	var reqErr error
	var appliance *openapi.Appliance
	if iok {
		appliance, reqErr = findApplianceByUUID(api, applianceID.(string), token)
	} else {
		appliance, reqErr = findApplianceByName(api, applianceName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got appliance: %+v", appliance)

	d.SetId(appliance.Id)
	d.Set("appliance_name", appliance.Name)
	d.Set("appliance_id", appliance.Id)
	return nil
}

func findApplianceByUUID(api *openapi.AppliancesApiService, id string, token string) (*openapi.Appliance, error) {
	log.Printf("[DEBUG] Data source appliance get by UUID %s", id)
	appliance, _, err := api.AppliancesIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &appliance, nil
}

func findApplianceByName(api *openapi.AppliancesApiService, name string, token string) (*openapi.Appliance, error) {
	log.Printf("[DEBUG] Data appliance get by name %s", name)
	request := api.AppliancesGet(context.Background())

	appliance, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range appliance.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find appliance %s", name)
}
