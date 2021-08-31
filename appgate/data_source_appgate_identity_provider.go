package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateIdentityProviderRead,
		Schema: map[string]*schema.Schema{
			"identity_provider_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"identity_provider_name"},
			},
			"identity_provider_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"identity_provider_id"},
			},
		},
	}
}

func dataSourceAppgateIdentityProviderRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source identity provider")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IdentityProvidersApi

	providerID, iok := d.GetOk("identity_provider_id")
	providerName, nok := d.GetOk("identity_provider_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of identity_provider_id or identity_provider_name attributes")
	}
	var reqErr error
	var provider *openapi.IdentityProvider
	if iok {
		provider, reqErr = findIdentityProviderByUUID(api, providerID.(string), token)
	} else {
		provider, reqErr = findIdentityProviderByName(api, providerName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got identity provider: %+v", provider)

	d.SetId(provider.Id)
	d.Set("identity_provider_name", provider.Name)
	d.Set("identity_provider_id", provider.Id)
	return nil
}

func findIdentityProviderByUUID(api *openapi.IdentityProvidersApiService, id string, token string) (*openapi.IdentityProvider, error) {
	provider, _, err := api.IdentityProvidersIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func findIdentityProviderByName(api *openapi.IdentityProvidersApiService, name string, token string) (*openapi.IdentityProvider, error) {
	request := api.IdentityProvidersGet(context.Background())

	provider, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range provider.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find identity provider %s", name)
}
