package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateMfaProvider() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateMfaProviderRead,
		Schema: map[string]*schema.Schema{
			"mfa_provider_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mfa_provider_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateMfaProviderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source MFA provider")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAProvidersApi

	providerID, iok := d.GetOk("mfa_provider_id")
	providerName, nok := d.GetOk("mfa_provider_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of mfa_provider_id or mfa_provider_name attributes")
	}
	var reqErr error
	var provider *openapi.MfaProvider
	if iok {
		provider, reqErr = findMfaProviderByUUID(api, providerID.(string), token)
	} else {
		provider, reqErr = findMfaProviderByName(api, providerName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got MFA provider: %+v", provider.Id)

	d.SetId(provider.Id)
	d.Set("mfa_provider_name", provider.Name)
	d.Set("mfa_provider_id", provider.Id)
	return nil
}

func findMfaProviderByUUID(api *openapi.MFAProvidersApiService, id string, token string) (*openapi.MfaProvider, error) {
	provider, _, err := api.MfaProvidersIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func findMfaProviderByName(api *openapi.MFAProvidersApiService, name string, token string) (*openapi.MfaProvider, error) {
	request := api.MfaProvidersGet(context.Background())

	provider, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range provider.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find MFA provider %s", name)
}
