package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v14/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateRingfenceRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateRingfenceRuleRead,
		Schema: map[string]*schema.Schema{

			"ringfence_rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ringfence_rule_name"},
			},

			"ringfence_rule_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ringfence_rule_id"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAppgateRingfenceRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source Ringfence Rules")
	token := meta.(*Client).Token
	api := meta.(*Client).API.RingfenceRulesApi
	ctx := context.Background()
	ringfenceID, iok := d.GetOk("ringfence_rule_id")
	ringfenceName, nok := d.GetOk("ringfence_rule_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of ringfence_rule_id or ringfence_rule_name attributes")
	}
	var reqErr error
	var ringfenceRule *openapi.RingfenceRule
	if iok {
		ringfenceRule, reqErr = findRingfenceRuleByUUID(ctx, api, ringfenceID.(string), token)
	} else {
		ringfenceRule, reqErr = findRingfenceRuleByName(ctx, api, ringfenceName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got Ringfence Rule: %+v", ringfenceRule)

	d.SetId(ringfenceRule.Id)
	d.Set("ringfence_rule_id", ringfenceRule.Id)
	d.Set("name", ringfenceRule.Name)
	d.Set("tags", ringfenceRule.Tags)

	return nil
}

func findRingfenceRuleByUUID(ctx context.Context, api *openapi.RingfenceRulesApiService, id string, token string) (*openapi.RingfenceRule, error) {
	log.Printf("[DEBUG] Data source Ringfence Rule get by UUID %s", id)
	ringfenceRule, _, err := api.RingfenceRulesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &ringfenceRule, nil
}

func findRingfenceRuleByName(ctx context.Context, api *openapi.RingfenceRulesApiService, name string, token string) (*openapi.RingfenceRule, error) {
	log.Printf("[DEBUG] Data source Ringfence Rule get by name %s", name)
	request := api.RingfenceRulesGet(ctx)

	ringfenceRule, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, r := range ringfenceRule.GetData() {
		return &r, nil
	}
	return nil, fmt.Errorf("Failed to find Ringfence rule %s", name)
}
