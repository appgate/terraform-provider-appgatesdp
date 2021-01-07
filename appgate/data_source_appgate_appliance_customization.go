package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateApplianceCustomization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateApplianceCustomizationRead,
		Schema: map[string]*schema.Schema{
			"appliance_customization_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"appliance_customization_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateApplianceCustomizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source Appliance customization")
	token := meta.(*Client).Token
	api := meta.(*Client).API.ApplianceCustomizationsApi

	applianceID, iok := d.GetOk("appliance_customization_id")
	applianceName, nok := d.GetOk("appliance_customization_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of appliance_customization_id or appliance_customization_name attributes")
	}
	var reqErr error
	var appliance *openapi.ApplianceCustomization
	if iok {
		appliance, reqErr = findApplianceCustomizationByUUID(api, applianceID.(string), token)
	} else {
		appliance, reqErr = findApplianceCustomizationByName(api, applianceName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got appliance customization: %+v", appliance)

	d.SetId(appliance.Id)
	d.Set("appliance_customization_name", appliance.Name)
	d.Set("appliance_customization_id", appliance.Id)
	return nil
}

func findApplianceCustomizationByUUID(api *openapi.ApplianceCustomizationsApiService, id string, token string) (*openapi.ApplianceCustomization, error) {
	log.Printf("[DEBUG] Data source appliance get by UUID %s", id)
	appliance, _, err := api.ApplianceCustomizationsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &appliance, nil
}

func findApplianceCustomizationByName(api *openapi.ApplianceCustomizationsApiService, name string, token string) (*openapi.ApplianceCustomization, error) {
	log.Printf("[DEBUG] Data appliance get by name %s", name)
	request := api.ApplianceCustomizationsGet(context.Background())

	appliance, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range appliance.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find appliance customization %s", name)
}
