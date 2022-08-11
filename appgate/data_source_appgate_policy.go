package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgatePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgatePolicyRead,
		Schema: map[string]*schema.Schema{

			"policy_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"policy_name"},
			},

			"policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"policy_id"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAppgatePolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source policy")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.PoliciesApi
	ctx := context.Background()

	policyID, iok := d.GetOk("policy_id")
	policyName, nok := d.GetOk("policy_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of policy_id or policy_name attributes")
	}

	var reqErr error
	var policy *openapi.Policy
	if iok {
		policy, reqErr = findPolicyByUUID(ctx, api, policyID.(string), token)
	} else {
		policy, reqErr = findPolicyByName(ctx, api, policyName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got policy: %+v", policy)

	d.SetId(policy.GetId())
	d.Set("name", policy.GetName())
	d.Set("policy_id", policy.GetId())
	d.Set("tags", policy.GetTags())

	return nil
}

func findPolicyByUUID(ctx context.Context, api *openapi.PoliciesApiService, id string, token string) (*openapi.Policy, error) {
	log.Printf("[DEBUG] Data source policy get by UUID %s", id)
	policy, _, err := api.PoliciesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func findPolicyByName(ctx context.Context, api *openapi.PoliciesApiService, name string, token string) (*openapi.Policy, error) {
	log.Printf("[DEBUG] Data source policy get by name %s", name)
	request := api.PoliciesGet(context.Background())

	policy, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, s := range policy.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find policy %s", name)
}
