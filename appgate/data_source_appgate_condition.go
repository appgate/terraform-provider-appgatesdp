package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAppgateCondition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateConditionRead,
		Schema: map[string]*schema.Schema{
			"condition_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"condition_name"},
			},
			"condition_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"condition_id"},
			},
		},
	}
}

func dataSourceAppgateConditionRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.ConditionsApi

	conditionID, iok := d.GetOk("condition_id")
	conditionName, nok := d.GetOk("condition_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of condition_id or condition_name attributes")
	}
	var reqErr error
	var condition *openapi.Condition
	if iok {
		condition, reqErr = findConditionByUUID(api, conditionID.(string), token)
	} else {
		condition, reqErr = findConditionByName(api, conditionName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(condition.Id)
	d.Set("name", condition.Name)

	return nil
}

func findConditionByUUID(api *openapi.ConditionsApiService, id string, token string) (*openapi.Condition, error) {
	condition, _, err := api.ConditionsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &condition, nil
}

func findConditionByName(api *openapi.ConditionsApiService, name string, token string) (*openapi.Condition, error) {
	request := api.ConditionsGet(context.Background())

	condition, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range condition.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find condition %s", name)
}
